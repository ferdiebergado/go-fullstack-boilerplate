import { clearErrors, showError } from "./form";
import { showNotification } from "./notification";
import { isRequiredInputFilled, isValidEmail } from "./validation";

const frm = document.getElementById("frmSignin") as HTMLFormElement;
const inputEmail = frm.querySelector("#email") as HTMLInputElement;
const inputPassword = frm.querySelector("#password") as HTMLInputElement;
const btnSignin = frm.querySelector("#btnSignin") as HTMLButtonElement;

let isLoading = false;

async function signInUser(e: SubmitEvent) {
	e.preventDefault();
	isLoading = true;
	updateSigninBtn();
	clearErrors(frm);

	const email = inputEmail.value;
	const password = inputPassword.value;

	const emailErr: string[] = [];
	const passwordErr: string[] = [];

	let isValid = true;

	if (!isRequiredInputFilled(inputEmail)) {
		emailErr.push("Email is required");
		isValid = false;
	}

	if (!isRequiredInputFilled(inputPassword)) {
		passwordErr.push("Password is required");
		isValid = false;
	}

	if (!isValidEmail(email)) {
		emailErr.push("Email must be a valid email address");
		isValid = false;
	}

	if (!isValid) {
		if (emailErr.length > 0) showError(inputEmail, emailErr);
		if (passwordErr.length > 0) showError(inputPassword, passwordErr);

		showNotification("error", "Invalid input!");
		isLoading = false;
		updateSigninBtn();
		return;
	}

	try {
		const res = await fetch(frm.action, {
			method: frm.method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email,
				password,
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
		updateSigninBtn();
	}
}

function updateSigninBtn(): void {
	btnSignin.disabled = isLoading;
	let btnText = "Sign In";
	if (isLoading) btnText = "Signing in...";

	btnSignin.textContent = btnText;
}

frm.addEventListener("submit", signInUser);

console.log("Waiting for user to sign in...");
