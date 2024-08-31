
export const showSpinner = () => document.getElementById('spinner').style.display = 'block';

export const hideSpinner = () => document.getElementById('spinner').style.display = 'none';

export const showPanelSpinner = panelElement => {
    const spinner = panelElement.querySelector('.panel-spinner');
    const imgElement = panelElement.querySelector('.grid-image');
    spinner.style.display = 'block';
    imgElement.style.display = 'none';
}

export const hidePanelSpinner = panelElement => {
    const spinner = panelElement.querySelector('.panel-spinner');
    const imgElement = panelElement.querySelector('.grid-image');
    spinner.style.display = 'none';
    imgElement.style.display = 'block';
}
