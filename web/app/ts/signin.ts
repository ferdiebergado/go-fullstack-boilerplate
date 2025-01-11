import { clearFormErrors, showFormError, updateSubmitBtn } from "./form";
import { showNotification } from "./notification";
import { isRequiredInputFilled, isValidEmail } from "./validation";

const frm = document.getElementById("frmSignin") as HTMLFormElement;
const inputEmail = frm.querySelector("#email") as HTMLInputElement;
const inputPassword = frm.querySelector("#password") as HTMLInputElement;
const btnSignin = frm.querySelector("#btnSignin") as HTMLButtonElement;

const btnAttrs = {
	btn: btnSignin,
	text: "Sign In",
	loadingText: "Signing in...",
};
const emailErr: string[] = [];
const passwordErr: string[] = [];

let isLoading = false;

function validate(): boolean {
	let isValid = true;

	if (!isRequiredInputFilled(inputEmail)) {
		emailErr.push("Email is required");
		isValid = false;
	}

	if (!isRequiredInputFilled(inputPassword)) {
		passwordErr.push("Password is required");
		isValid = false;
	}

	if (!isValidEmail(inputEmail.value)) {
		emailErr.push("Email must be a valid email address");
		isValid = false;
	}

	return isValid;
}

async function signInUser(e: SubmitEvent) {
	e.preventDefault();
	isLoading = true;
	updateSubmitBtn(btnAttrs, isLoading);
	clearFormErrors(frm);

	const isValid = validate();

	if (!isValid) {
		if (emailErr.length > 0) showFormError(inputEmail, emailErr);
		if (passwordErr.length > 0) showFormError(inputPassword, passwordErr);

		showNotification("error", "Invalid input!");
		isLoading = false;
		updateSubmitBtn(btnAttrs, isLoading);
		return;
	}

	try {
		const res = await fetch(frm.action, {
			method: frm.method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: inputEmail.value,
				password: inputPassword.value,
			}),
		});

		if (!res.ok) {
			const { message, errors }: APIResponse<undefined> =
				await res.json();

			if (errors) {
				for (const field in errors) {
					console.log("field", field, "errors", errors[field]);
					const el = document.getElementById(field);
					if (el) {
						showFormError(
							document.getElementById(field) as HTMLInputElement,
							errors[field]
						);
					}
				}
			}

			showNotification("error", message);
		} else {
			const { message, data }: APIResponse<undefined> = await res.json();

			console.log("message", message, "data", data);

			showNotification("success", message);
		}
	} catch (error) {
		console.log(error);
		if (error instanceof Error) showNotification("error", error.message);
	} finally {
		isLoading = false;
		updateSubmitBtn(btnAttrs, isLoading);
	}
}

frm.addEventListener("submit", signInUser);

console.log("Waiting for user to sign in...");
