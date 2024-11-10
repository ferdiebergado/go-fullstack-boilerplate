// @ts-check
'use strict';

const footer = document.querySelector('footer');
const year = new Date().getFullYear();

footer && (footer.innerHTML = `<span>&copy; ${year}</span>`);
