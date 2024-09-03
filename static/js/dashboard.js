import { parsePanel } from "./panel";

export const DASHBOARD = window.location.pathname.split('/').slice(-2, -1)[0];
export const SCREEN = window.innerWidth;

const ROW_NUM = 24;  // Grafana's max x
const OPTIONS = {
    float: false,
    cellHeight: 38,
    margin: 4,
    column: ROW_NUM,
}

export const getDashboard = async () => {
    const url = `/api/v1/dashboard/${DASHBOARD}/?screen=${SCREEN}`;
    const response = await fetch(url, { credentials: 'include' });
    if (!response.ok) throw new Error(`Can't load dashboard: ${response.status}`);
    return await response.json();
}

export const loadDashboard = (data) => {
    const items = data.panels.map(parsePanel);
    window.grid = GridStack.init(OPTIONS).load(items); 
};