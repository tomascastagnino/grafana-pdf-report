
const ROW_NUM = 12;

document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const dashboard_uid = window.location.pathname.split('/').slice(-2, -1)[0];
    const params = urlParams.toString();

    const apiUrl = `/api/v1/report/data/${dashboard_uid}/?${params}`;

    const spinner = document.getElementById('spinner');

    try {
        spinner.style.display = 'block';

        const response = await fetch(apiUrl);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        const options = {
            float: false,
            cellHeight: 50,
            margin: 4,
            column: ROW_NUM,
            padding: 20
        };
        const items = [];

        for (const panel of Object.values(data.panels)) {
            const imageUrl = panel.url;

            let button = `<a href="#" class="close-button" onclick="removePanel(this)"></a>`;
            let movement = {};

            if (panel.tag === "fixed") {
                movement = {
                    noMove: true,
                    noResize: true,
                    locked: true,
                }
                button = `<div></div>`
            }
            const innerContent = panel.type === "text" ?
                                        `<div class="text-panel">${panel.options.content}</div>` :
                                        `<img src="${imageUrl}" class="grid-image">`

            const content = `
                <div class="image-container" data-panel-id="${panel.id}">
                    ${button}
                    ${innerContent}
                </div>`;

            let panelObj = {
                x: parseInt(panel.gridPos.x) / 2.0,
                y: parseInt(panel.gridPos.y),
                w: parseInt(panel.gridPos.w) / 2.0,
                h: parseInt(panel.gridPos.h),
                ...movement,
                content: content
            }
            items.push(panelObj);
        }

        window.grid = GridStack.init(options).load(items); 

        window.removePanel = function(button) {
            const panelElement = button.closest('.grid-stack-item');
            if (panelElement) {
                window.grid.removeWidget(panelElement);
            }
        };

        spinner.style.display = 'none';

        const generatePdfButton = document.getElementById('generatePdfButton');
        generatePdfButton.addEventListener('click', async () => {
            const gridElement = document.querySelector('.grid-stack');
            const images = gridElement.querySelectorAll('.grid-image');
            const closeButtons = gridElement.querySelectorAll('.close-button');

            // Hide the close buttons
            closeButtons.forEach(button => button.style.display = 'none');

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

            const gridWidth = gridElement.offsetWidth;
            const gridHeight = gridElement.offsetHeight;

            const opt = {
                margin: 1,
                filename: 'dashboard.pdf',
                image: { type: 'jpeg', quality: 0.98 },
                html2canvas: { scale: 2, useCORS: true },
                jsPDF: { unit: 'px', format: [gridWidth, gridHeight], orientation: 'portrait' }
            };

            await html2pdf().set(opt).from(gridElement).save();
            // Restore the close buttons
            closeButtons.forEach(button => button.style.display = 'inline-block');
        });
    } catch (error) {
        spinner.style.display = 'none'; 
        console.error('Error fetching dashboard data:', error);
    }
});