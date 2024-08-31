
export async function createPDF() {
    const gridElement = document.querySelector('.grid-stack');
    const images = gridElement.querySelectorAll('.grid-image');
    const closeButtons = gridElement.querySelectorAll('.close-button');
    const refreshButtons = gridElement.querySelectorAll('.refresh-button');

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

    closeButtons.forEach(button => button.style.display = 'inline-block');
    refreshButtons.forEach(button => button.style.display = 'inline-block');
}
