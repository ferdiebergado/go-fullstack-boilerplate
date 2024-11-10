// @ts-check
'use strict';

const footer = document.getElementsByTagName('footer')[0];

const year = new Date().getFullYear();

footer.innerHTML = `&copy; ${year}`;
