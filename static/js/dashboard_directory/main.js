document.addEventListener('DOMContentLoaded', async () => {
    const response = await fetch('/api/v1/dashboards');
    const items = await response.json();

    const list = document.getElementById('dashboardList');

    const idMap = {};
    const root = [];

    items.forEach(item => {
        idMap[item.uid] = item;
        item.content = [];
    });

    // I'm creating a tree structure by adding one object inside another.
    // Then, I place the root objects in an array.
    items.forEach(item => {
        if (!item.folderUid) {
            root.push(item);
        } else {
            idMap[item.folderUid].content.push(item);
        }
    });

    const buildFolder = (item, el) => {
        el.classList.add('folder');
        el.textContent = item.title.toLowerCase().replaceAll(' ', '_') + '/';
    };

    const buildFile = (item, el) => {
        el.classList.add('dashboard');
        el.textContent = item.title.toLowerCase().replaceAll(' ', '_');
        el.href = `/api/v1/report/${item.uid}/`;
    };

    // Recursive function. If the object is not-a-file then break, else continue appending 'li' elements.
    const build = (item, dom) => {
        const el = document.createElement('li');
        if (item.type !== 'dash-folder') {
            const link = document.createElement('a');
            buildFile(item, link);
            el.appendChild(link);
        } else {
            buildFolder(item, el);
            const branch = document.createElement('ul');
            el.appendChild(branch);
            item.content.forEach(subItem => build(subItem, branch));
        }
        dom.appendChild(el);
    };

    root.forEach(item => {
        build(item, list);
    });
});
