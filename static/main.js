document.addEventListener("DOMContentLoaded", function() {
    const panelsDiv = document.getElementById('panels');
    const generatePdfButton = document.getElementById('generate-pdf');

    const urlParams = new URLSearchParams(window.location.search);
    const pathParts = window.location.pathname.split('/');
    const dashboardID = pathParts[4];  // TODO: refactor this 
    const queryParams = urlParams.toString();

    fetch(`/api/v1/report/data/${dashboardID}?${queryParams}`)
        .then(response => response.json())
        .then(data => {
            for (const [panelID, imageURL] of Object.entries(data.panel_images)) {
                const img = document.createElement('img');
                img.src = imageURL;
                img.style.width = '300px';
                img.style.height = '150px';

                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.id = `panel-${panelID}`;
                checkbox.value = imageURL;

                const label = document.createElement('label');
                label.htmlFor = `panel-${panelID}`;
                label.appendChild(img);

                const div = document.createElement('div');
                div.appendChild(checkbox);
                div.appendChild(label);

                panelsDiv.appendChild(div);
            }

            generatePdfButton.style.display = 'block';
        });

    generatePdfButton.addEventListener('click', function() {
        const selectedPanels = [];
        document.querySelectorAll('input[type="checkbox"]:checked').forEach(checkbox => {
            selectedPanels.push(checkbox.value);
        });

        fetch('/generate-pdf', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ panels: selectedPanels })
        })
        .then(response => response.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'selected_panels.pdf';
            document.body.appendChild(a);
            a.click();
            a.remove();
        });
    });
});