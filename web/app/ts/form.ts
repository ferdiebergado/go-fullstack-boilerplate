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

export function handleFormErrors(
	errors: ValidationErrorMap,
	form: HTMLFormElement
) {
	for (const field in errors) {
		if (errors[field].length > 0) {
			const input = form.querySelector(`#${field}`) as HTMLInputElement;
			if (input) {
				showFormError(input, errors[field]);
			}
		}
	}
}

export function toggleError(input: HTMLInputElement, message: string) {
	const helpText = input.nextElementSibling as HTMLElement;
	if (message) {
		input.classList.add("error");
		helpText.textContent = message;
		helpText.style.display = "block";
	} else {
		input.classList.remove("error");
		helpText.textContent = "";
		helpText.style.display = "none";
	}
}
