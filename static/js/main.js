document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const dashboard_uid = window.location.pathname.split('/').slice(-2, -1)[0];
    const params = urlParams.toString();

    const apiUrl = `/api/v1/report/data/${dashboard_uid}/?${params}`;
    
    try {
        const response = await fetch(apiUrl);
        const data = await response.json();
        
        const gridContainer = document.getElementById('gridContainer');
        
        for (const panel of Object.values(data.panels)) {
            const panelDiv = document.createElement('div');
            panelDiv.style.gridColumnStart = panel.gridPos.x + 1;
            panelDiv.style.gridColumnEnd = `span ${panel.gridPos.w}`;
            panelDiv.style.gridRowStart = panel.gridPos.y + 1;
            panelDiv.style.gridRowEnd = `span ${panel.gridPos.h}`;
            const contentDiv = document.createElement('div'); 
            const checkbox = document.createElement('input');
            panelDiv.appendChild(checkbox);
            panelDiv.className = 'grid-item';
            if (panel.tag !== "fixed") {
                checkbox.type = 'checkbox';
                checkbox.className = 'checkbox';
                checkbox.checked = true;
                checkbox.value = panel.id; 
            }
            if (panel.type === "text") {
                panelDiv.classList.add('text-panel');
                panelDiv.style.height = `calc(${panel.gridPos.h} * 37.11px)`
                contentDiv.innerHTML = panel.options.content;
            } else {
                const img = document.createElement('img');
                img.src = panel.url;
                panelDiv.appendChild(img);
            } 
            panelDiv.appendChild(contentDiv);
            gridContainer.appendChild(panelDiv);
        }

        const generatePdfButton = document.getElementById('generatePdfButton');
        generatePdfButton.addEventListener('click', async () => {
            const selectedPanels = [];
            document.querySelectorAll('input[type="checkbox"]:checked').forEach(checkbox => {
                const panelID = checkbox.value;
                const panel = data.panels[panelID];
                selectedPanels.push({
                    url: panel.url,
                    gridPos: panel.gridPos,
                    id: Number(panelID)
                });
            });

            try {
                const response = await fetch('/generate-pdf', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ panels: selectedPanels })
                });

                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }

                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'selected_panels.pdf';
                document.body.appendChild(a);
                a.click();
                a.remove();
            } catch (error) {
                console.error('There was a problem with the fetch operation:', error);
            }
        });
    } catch (error) {
        console.error('Error fetching dashboard data:', error);
    }
});
