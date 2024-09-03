import { showPanelSpinner, hidePanelSpinner } from './spinner.js';
import { DASHBOARD, SCREEN } from './dashboard.js';



export const refreshPanel = async (button) => {
    const panelElement = button.closest('.grid-stack-item');
    const panel = panelElement.gridstackNode;

    const params = new URLSearchParams(window.location.search);

    params.append("screen", SCREEN);
    params.append("w", panel.w);
    params.append("h", panel.h);

    showPanelSpinner(panelElement);

    const data = await getPanel(panel.id, params);

    panel.url = data.url;
 
    window.grid.update(panelElement, {"content": getContent(panel)})

    hidePanelSpinner(panelElement);
};

export const removePanel = (button) => {
    const panelElement = button.closest('.grid-stack-item');
    if (panelElement) {
        window.grid.removeWidget(panelElement);
    }
};

const getPanel = async (panel, params) => {
    const url = `/api/v1/dashboard/${DASHBOARD}/panel/${panel}/?${params}`;
    const response = await fetch(url, { credentials: 'include' });
    if (!response.ok) throw new Error(`Failed to refresh panel: ${response.status}`);
    return await response.json();
}

const getContent = (panel) => {
    const closeButton = '<button href="#" class="close-button" onclick="removePanel(this)">x</button>';
    const refreshButton = '<button href="#" class="refresh-button" onclick="refreshPanel(this)">&#x21bb;</button>';
    const spinner = '<div class="panel-spinner"><div class="panel-spinner-circle"></div></div>';
    const imageUrl = panel.url;

    return `
        <div class="image-container" data-panel-id="${panel.id}">
            ${refreshButton}
            ${closeButton}
            <img src="${imageUrl}" class="grid-image">${spinner}
        </div>`;
}

export function parsePanel(panel) {
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