export function showError(el: HTMLInputElement, errors: string[]): void {
	el.classList.add("error");

	const helpText = el.nextElementSibling as HTMLElement;

	if (helpText) {
		Array.from({ length: errors.length }, (_, i) => {
			helpText.textContent += errors[i] + "\n";
		});
		helpText.style.display = "block";
	}
}

export function clearErrors(frm: HTMLFormElement): void {
	frm.querySelectorAll(".error").forEach((e) => {
		e.classList.remove("error");
	});

	frm.querySelectorAll<HTMLElement>(".help-text").forEach((e) => {
		e.textContent = "";
		e.style.display = "none";
	});
}
