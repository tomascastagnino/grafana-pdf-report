import { getDashboard, loadDashboard } from './dashboard.js';
import { refreshPanel, removePanel } from './panel.js';
import { showSpinner, hideSpinner } from './spinner.js';
import { createPDF } from './pdf.js';


document.addEventListener('DOMContentLoaded', async () => {
    try {
        showSpinner();

        const data = await getDashboard();
        loadDashboard(data);

        hideSpinner();

        const generatePdfButton = document.getElementById('generatePdfButton');
        generatePdfButton.addEventListener('click', createPDF);

    } catch (error) {
        hideSpinner();
        console.error('Error fetching dashboard data:', error);
    }
});

window.removePanel = removePanel;
window.refreshPanel = refreshPanel; 
