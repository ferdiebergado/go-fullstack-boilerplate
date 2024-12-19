const spanYear = document.getElementById("year");

export function getCurrentYear(): number {
	return new Date().getFullYear();
}

spanYear && (spanYear.textContent = getCurrentYear().toString());
