/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./static/js/dashboard.js":
/*!********************************!*\
  !*** ./static/js/dashboard.js ***!
  \********************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   DASHBOARD: () => (/* binding */ DASHBOARD),\n/* harmony export */   SCREEN: () => (/* binding */ SCREEN),\n/* harmony export */   getDashboard: () => (/* binding */ getDashboard),\n/* harmony export */   loadDashboard: () => (/* binding */ loadDashboard)\n/* harmony export */ });\n/* harmony import */ var _panel__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./panel */ \"./static/js/panel.js\");\n\n\nconst DASHBOARD = window.location.pathname.split('/').slice(-2, -1)[0];\nconst SCREEN = window.innerWidth;\n\nconst ROW_NUM = 24;  // Grafana's max x\nconst OPTIONS = {\n    float: false,\n    cellHeight: 38,\n    margin: 4,\n    column: ROW_NUM,\n}\n\nconst getDashboard = async () => {\n    const url = `/api/v1/dashboard/${DASHBOARD}/?screen=${SCREEN}`;\n    const response = await fetch(url, { credentials: 'include' });\n    if (!response.ok) throw new Error(`Can't load dashboard: ${response.status}`);\n    return await response.json();\n}\n\nconst loadDashboard = (data) => {\n    const items = data.panels.map(_panel__WEBPACK_IMPORTED_MODULE_0__.parsePanel);\n    window.grid = GridStack.init(OPTIONS).load(items); \n};\n\n//# sourceURL=webpack://grafana-pdf-reporter/./static/js/dashboard.js?");

/***/ }),

/***/ "./static/js/main.js":
/*!***************************!*\
  !*** ./static/js/main.js ***!
  \***************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _dashboard_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./dashboard.js */ \"./static/js/dashboard.js\");\n/* harmony import */ var _panel_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./panel.js */ \"./static/js/panel.js\");\n/* harmony import */ var _spinner_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./spinner.js */ \"./static/js/spinner.js\");\n/* harmony import */ var _pdf_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./pdf.js */ \"./static/js/pdf.js\");\n\n\n\n\n\n\ndocument.addEventListener('DOMContentLoaded', async () => {\n    try {\n        (0,_spinner_js__WEBPACK_IMPORTED_MODULE_2__.showSpinner)();\n\n        const data = await (0,_dashboard_js__WEBPACK_IMPORTED_MODULE_0__.getDashboard)();\n        (0,_dashboard_js__WEBPACK_IMPORTED_MODULE_0__.loadDashboard)(data);\n\n        (0,_spinner_js__WEBPACK_IMPORTED_MODULE_2__.hideSpinner)();\n\n        const generatePdfButton = document.getElementById('generatePdfButton');\n        generatePdfButton.addEventListener('click', _pdf_js__WEBPACK_IMPORTED_MODULE_3__.createPDF);\n\n    } catch (error) {\n        (0,_spinner_js__WEBPACK_IMPORTED_MODULE_2__.hideSpinner)();\n        console.error('Error fetching dashboard data:', error);\n    }\n});\n\nwindow.removePanel = _panel_js__WEBPACK_IMPORTED_MODULE_1__.removePanel;\nwindow.refreshPanel = _panel_js__WEBPACK_IMPORTED_MODULE_1__.refreshPanel; \n\n\n//# sourceURL=webpack://grafana-pdf-reporter/./static/js/main.js?");

/***/ }),

/***/ "./static/js/panel.js":
/*!****************************!*\
  !*** ./static/js/panel.js ***!
  \****************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   parsePanel: () => (/* binding */ parsePanel),\n/* harmony export */   refreshPanel: () => (/* binding */ refreshPanel),\n/* harmony export */   removePanel: () => (/* binding */ removePanel)\n/* harmony export */ });\n/* harmony import */ var _spinner_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./spinner.js */ \"./static/js/spinner.js\");\n/* harmony import */ var _dashboard_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./dashboard.js */ \"./static/js/dashboard.js\");\n\n\n\n\n\nconst refreshPanel = async (button) => {\n    const panelElement = button.closest('.grid-stack-item');\n    const panel = panelElement.gridstackNode;\n\n    const params = new URLSearchParams(window.location.search);\n\n    params.append(\"screen\", _dashboard_js__WEBPACK_IMPORTED_MODULE_1__.SCREEN);\n    params.append(\"w\", panel.w);\n    params.append(\"h\", panel.h);\n\n    (0,_spinner_js__WEBPACK_IMPORTED_MODULE_0__.showPanelSpinner)(panelElement);\n\n    const data = await getPanel(panel.id, params);\n\n    panel.url = data.url;\n \n    window.grid.update(panelElement, {\"content\": getContent(panel)})\n\n    ;(0,_spinner_js__WEBPACK_IMPORTED_MODULE_0__.hidePanelSpinner)(panelElement);\n};\n\nconst removePanel = (button) => {\n    const panelElement = button.closest('.grid-stack-item');\n    if (panelElement) {\n        window.grid.removeWidget(panelElement);\n    }\n};\n\nconst getPanel = async (panel, params) => {\n    const url = `/api/v1/dashboard/${_dashboard_js__WEBPACK_IMPORTED_MODULE_1__.DASHBOARD}/panel/${panel}/?${params}`;\n    const response = await fetch(url, { credentials: 'include' });\n    if (!response.ok) throw new Error(`Failed to refresh panel: ${response.status}`);\n    return await response.json();\n}\n\nconst getContent = (panel) => {\n    const closeButton = '<button href=\"#\" class=\"close-button\" onclick=\"removePanel(this)\">x</button>';\n    const refreshButton = '<button href=\"#\" class=\"refresh-button\" onclick=\"refreshPanel(this)\">&#x21bb;</button>';\n    const spinner = '<div class=\"panel-spinner\"><div class=\"panel-spinner-circle\"></div></div>';\n    const imageUrl = panel.url;\n\n    return `\n        <div class=\"image-container\" data-panel-id=\"${panel.id}\">\n            ${refreshButton}\n            ${closeButton}\n            <img src=\"${imageUrl}\" class=\"grid-image\">${spinner}\n        </div>`;\n}\n\nfunction parsePanel(panel) {\n    let movement = {};\n    if (panel.tag === \"fixed\") {\n        movement = {\n            noMove: true,\n            noResize: true,\n            locked: true,\n        }\n    }\n    return {\n        x: parseInt(panel.gridPos.x),\n        y: parseInt(panel.gridPos.y),\n        w: parseInt(panel.gridPos.w),\n        h: parseInt(panel.gridPos.h),\n        id: panel.id,\n        tag: panel.tag,\n        type: panel.type,\n        ...movement,\n        content: getContent(panel)\n    }\n}\n\n//# sourceURL=webpack://grafana-pdf-reporter/./static/js/panel.js?");

/***/ }),

/***/ "./static/js/pdf.js":
/*!**************************!*\
  !*** ./static/js/pdf.js ***!
  \**************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   createPDF: () => (/* binding */ createPDF)\n/* harmony export */ });\n\nasync function createPDF() {\n    const gridElement = document.querySelector('.grid-stack');\n    const images = gridElement.querySelectorAll('.grid-image');\n    const closeButtons = gridElement.querySelectorAll('.close-button');\n    const refreshButtons = gridElement.querySelectorAll('.refresh-button');\n\n    closeButtons.forEach(button => button.style.display = 'none');\n    refreshButtons.forEach(button => button.style.display = 'none');\n\n    await Promise.all(Array.from(images).map(img => {\n        return new Promise(resolve => {\n            if (img.complete) {\n                resolve();\n            } else {\n                img.onload = resolve;\n                img.onerror = resolve;\n            }\n        });\n    }));\n\n    const gridWidth = gridElement.scrollWidth;\n    const gridHeight = gridElement.scrollHeight;\n\n    gridElement.style.width = `${gridWidth}px`;\n    gridElement.style.height = `${gridHeight}px`;\n\n    const opt = {\n        margin: 20,\n        filename: 'dashboard.pdf',\n        image: { type: 'jpeg', quality: 1 },\n        html2canvas: { scale: 2, width: gridWidth, height: gridHeight, logging: false },\n        jsPDF: { unit: 'px', format: [gridWidth, gridHeight], orientation: 'portrait' }\n    };\n\n    await html2pdf().set(opt).from(gridElement).save();\n\n    closeButtons.forEach(button => button.style.display = 'inline-block');\n    refreshButtons.forEach(button => button.style.display = 'inline-block');\n}\n\n\n//# sourceURL=webpack://grafana-pdf-reporter/./static/js/pdf.js?");

/***/ }),

/***/ "./static/js/spinner.js":
/*!******************************!*\
  !*** ./static/js/spinner.js ***!
  \******************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   hidePanelSpinner: () => (/* binding */ hidePanelSpinner),\n/* harmony export */   hideSpinner: () => (/* binding */ hideSpinner),\n/* harmony export */   showPanelSpinner: () => (/* binding */ showPanelSpinner),\n/* harmony export */   showSpinner: () => (/* binding */ showSpinner)\n/* harmony export */ });\n\nconst showSpinner = () => document.getElementById('spinner').style.display = 'block';\n\nconst hideSpinner = () => document.getElementById('spinner').style.display = 'none';\n\nconst showPanelSpinner = panelElement => {\n    const spinner = panelElement.querySelector('.panel-spinner');\n    const imgElement = panelElement.querySelector('.grid-image');\n    spinner.style.display = 'block';\n    imgElement.style.display = 'none';\n}\n\nconst hidePanelSpinner = panelElement => {\n    const spinner = panelElement.querySelector('.panel-spinner');\n    const imgElement = panelElement.querySelector('.grid-image');\n    spinner.style.display = 'none';\n    imgElement.style.display = 'block';\n}\n\n\n//# sourceURL=webpack://grafana-pdf-reporter/./static/js/spinner.js?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = __webpack_require__("./static/js/main.js");
/******/ 	
/******/ })()
;