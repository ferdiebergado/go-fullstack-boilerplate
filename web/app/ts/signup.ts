import { showNotification } from "./notification";

const frm = document.forms.namedItem("frmSignup") as HTMLFormElement;
const inputEmail = frm.querySelector("#email") as HTMLInputElement;
const inputPassword = frm.querySelector("#password") as HTMLInputElement;
const inputRetypePass = frm.querySelector(
	"#password_confirmation"
) as HTMLInputElement;
const btnSignup = frm.querySelector("#btnSignup") as HTMLButtonElement;

let isLoading = false;

async function signUpUser(e: SubmitEvent) {
	e.preventDefault();
	isLoading = true;
	updateSignupBtn();
	clearErrors();

	try {
		const res = await fetch(frm.action, {
			method: frm.method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: inputEmail.value,
				password: inputPassword.value,
				password_confirmation: inputRetypePass.value,
			}),
		});

		const {
			success,
			message,
			data,
		}: APIResponse<ValidationErrors | undefined> = await res.json();

		console.log("success", success, "message", message, "data", data);

		let notifType: NotificationType = "success";

		if (!success) {
			notifType = "error";
			if (data) {
				const { errors } = data;

				if (errors) {
					for (const field in errors) {
						console.log("field", field, "errors", errors[field]);
						const el = document.getElementById(field);
						if (el) {
							showError(
								document.getElementById(
									field
								) as HTMLInputElement,
								errors[field]
							);
						}
					}
				}
			}
		}

		showNotification(notifType, message);
	} catch (error) {
		console.log(error);
		if (error instanceof Error) showNotification("error", error.message);
	} finally {
		isLoading = false;
		updateSignupBtn();
	}
}

function comparePasswords(): void {
	if (inputRetypePass.value && inputRetypePass.value != inputPassword.value) {
		inputRetypePass.classList.add("error");
	} else {
		inputRetypePass.classList.remove("error");
	}
}

function updateSignupBtn(): void {
	btnSignup.disabled = isLoading;
	let btnText = "Sign Up";
	if (isLoading) btnText = "Signing up...";

	btnSignup.textContent = btnText;
}

function showError(el: HTMLInputElement, errors: string[]): void {
	el.classList.add("error");

	const helpText = el.nextElementSibling as HTMLElement;

	if (helpText) {
		Array.from({ length: errors.length }, (_, i) => {
			helpText.textContent += errors[i] + "\n";
		});
		helpText.style.display = "block";
	}
}

function clearErrors(): void {
	frm.querySelectorAll(".error").forEach((e) => {
		e.classList.remove("error");
	});

	frm.querySelectorAll<HTMLElement>(".help-text").forEach((e) => {
		e.textContent = "";
		e.style.display = "none";
	});
}

inputPassword.addEventListener("input", comparePasswords);
inputRetypePass.addEventListener("input", comparePasswords);
frm.addEventListener("submit", signUpUser);

console.log("Waiting for user to sign up...");
