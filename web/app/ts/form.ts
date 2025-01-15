import { resetValidationErrors } from "./validation";

export function showFormError(el: HTMLInputElement, errors: string[]): void {
	el.classList.add("error");

	const helpText = el.nextElementSibling as HTMLElement;

	if (helpText) {
		helpText.innerHTML = "";
		const fragment = document.createDocumentFragment();
		Array.from({ length: errors.length }, (_, i) => {
			const li = document.createElement("li");
			const error = errors[i];
			console.log("error:", error);
			li.textContent = error;
			fragment.appendChild(li);
		});
		resetValidationErrors(errors);
		helpText.appendChild(fragment);
		helpText.style.display = "block";
	}
}

export function clearFormErrors(frm: HTMLFormElement): void {
	frm.querySelectorAll(".error").forEach((e) => {
		e.classList.remove("error");
	});

	frm.querySelectorAll<HTMLElement>(".help-text").forEach((e) => {
		e.textContent = "";
		e.style.display = "none";
	});
}

export function updateSubmitBtn(
	{ btn, text, loadingText }: ButtonAttrs,
	isLoading: boolean
): void {
	btn.disabled = isLoading;
	let btnText = text;
	if (isLoading) btnText = loadingText;

	btn.textContent = btnText;
}
