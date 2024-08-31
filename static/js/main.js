
const ROW_NUM = 24;  // Grafana's max x
const screen = window.innerWidth;
const urlParams = new URLSearchParams(window.location.search);
let dashboardId;

const getContent = (panel) => {
    const closeButton = '<button href="#" class="close-button" onclick="removePanel(this)">x</button>';
    const refreshButton = '<button href="#" class="refresh-button" onclick="refreshPanel(this)">&#x21bb;</button>';
    const spinner = '<div class="panel-spinner"><div class="panel-spinner-circle"></div></div>';
    const imageUrl = panel.url;

    // I need to sanitize this
    // const innerContent = panel.type === "text" ?
    //     `<div class="text-panel">${panel.options.content}</div>` :
    //     `<img src="${imageUrl}" class="grid-image">${spinner}`

    const innerContent = `<img src="${imageUrl}" class="grid-image">${spinner}`

    return `
        <div class="image-container" data-panel-id="${panel.id}">
            ${refreshButton}
            ${closeButton}
            ${innerContent}
        </div>`;
}

const parsedPanel = (panel) => {
    let movement = {};
    if (panel.tag === "fixed") {
        movement = {
            noMove: true,
            noResize: true,
            locked: true,
        }
    }
    return {
        x: parseInt(panel.gridPos.x),
        y: parseInt(panel.gridPos.y),
        w: parseInt(panel.gridPos.w),
        h: parseInt(panel.gridPos.h),
        id: panel.id,
        tag: panel.tag,
        type: panel.type,
        ...movement,
        content: getContent(panel)
    }
}

document.addEventListener('DOMContentLoaded', async () => {
    dashboardId = window.location.pathname.split('/').slice(-2, -1)[0];
    urlParams.append("screen", screen);
    const params = urlParams.toString();

    const apiUrl = `/api/v1/dashboard/${dashboardId}/?${params}`;

    const spinner = document.getElementById('spinner');

    try {
        spinner.style.display = 'block';

        const response = await fetch(apiUrl, {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();

        const options = {
            float: false,
            cellHeight: 38, // Grafana's aprox height to widht ratio
            margin: 4,
            column: ROW_NUM,
        };

        const items = [];
        for (const panel of Object.values(data.panels)) {
            items.push(parsedPanel(panel));
        }

        window.grid = GridStack.init(options).load(items); 

        spinner.style.display = 'none';

        const generatePdfButton = document.getElementById('generatePdfButton');
        generatePdfButton.addEventListener('click', createPDF)

    } catch (error) {
        spinner.style.display = 'none'; 
        console.error('Error fetching dashboard data:', error);
    }
});

window.removePanel = (button) => {
    const panelElement = button.closest('.grid-stack-item');
    if (panelElement) {
        window.grid.removeWidget(panelElement);
    }
};

window.refreshPanel = async (button) => {
    const panelElement = button.closest('.grid-stack-item');
    const panel = panelElement.gridstackNode;

    let params = new URLSearchParams(Object.fromEntries(urlParams));
    params.append("screen", screen);

    if (panel) {
        const w = panel.w;
        const h = panel.h;
        params.append("w", w);
        params.append("h", h);
        params = params.toString();

        const apiUrl = `/api/v1/dashboard/${dashboardId}/panel/${panel.id}/?${params}`;

        const spinner = panelElement.querySelector('.panel-spinner');
        const imgElement = panelElement.querySelector('.grid-image');
        spinner.style.display = 'block';
        imgElement.style.display = 'none';

        try {
            const response = await fetch(apiUrl, {
                credentials: 'include'
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            panel.url = data.url;
            const content = getContent(panel);
         
            window.grid.update(panelElement, {"content": content})

            spinner.style.display = 'none';
            imgElement.style.display = 'block';

        } catch (error) {
            console.error('Error refreshing panel:', error);
        }
    }
};

const createPDF = async () => {
    const gridElement = document.querySelector('.grid-stack');
    const images = gridElement.querySelectorAll('.grid-image');
    const closeButtons = gridElement.querySelectorAll('.close-button');
    const refreshButtons = gridElement.querySelectorAll('.refresh-button');

    // Hide the buttons
    closeButtons.forEach(button => button.style.display = 'none');
    refreshButtons.forEach(button => button.style.display = 'none');

    await Promise.all(Array.from(images).map(img => {
        return new Promise(resolve => {
            if (img.complete) {
                resolve();
            } else {
                img.onload = resolve;
                img.onerror = resolve;
            }
        });
    }));

    const gridWidth = gridElement.scrollWidth;
    const gridHeight = gridElement.scrollHeight;

    // Set the width and height to get the correct rendering ratio
    gridElement.style.width = `${gridWidth}px`;
    gridElement.style.height = `${gridHeight}px`;

    const opt = {
        margin: 20,
        filename: 'dashboard.pdf',
        image: { type: 'jpeg', quality: 1 },
        html2canvas: { scale: 2, width: gridWidth, height: gridHeight, logging: false },
        jsPDF: { unit: 'px', format: [gridWidth, gridHeight], orientation: 'portrait' }
    };

    await html2pdf().set(opt).from(gridElement).save();

    // Restore the buttons
    closeButtons.forEach(button => button.style.display = 'inline-block');
    refreshButtons.forEach(button => button.style.display = 'inline-block');
};
