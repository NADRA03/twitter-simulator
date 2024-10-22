import * as auth from './auth.js';

async function handleRoute() {
    const path = window.location.pathname;

    if (path === "/") {
        auth.renderFirstPage();
    } else if (path === "/login") {
        auth.renderLoginTemplate();
    } else if (path === "/signup") {
        auth.renderSignInTemplate();
    } else if (path === "/home") {
        auth.renderHomePage();
    }
};

handleRoute();