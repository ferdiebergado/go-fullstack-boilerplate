import {
	clearFormErrors,
	handleFormErrors,
	showFormError,
	toggleError,
	updateSubmitBtn,
} from "./form";
import { showNotification } from "./notification";
import { isRequiredInputFilled, isValidEmail } from "./validation";

const frmSignin = document.getElementById("frmSignin") as HTMLFormElement;
const inputEmail = frmSignin.querySelector("#email") as HTMLInputElement;
const inputPassword = frmSignin.querySelector("#password") as HTMLInputElement;
const btnSignin = frmSignin.querySelector("#btnSignin") as HTMLButtonElement;

const btnAttrs = {
	btn: btnSignin,
	text: "Sign In",
	loadingText: "Signing in...",
};
type ValidationErrors = {
	email: Errors;
	password: Errors;
};

const validationErrors: ValidationErrors = {
	email: [],
	password: [],
};

let isLoading = false;

frmSignin.addEventListener("change", handleInputChange);
frmSignin.addEventListener("submit", signInUser);

function handleInputChange(event: Event) {
	const target = event.target as HTMLInputElement;
	if (target.matches("#email")) {
		toggleError(event.target as HTMLInputElement, "");

		if (!isValidEmail(target.value)) {
			showFormError(target, ["Email must be a valid email address"]);
		}
	}
}

async function signInUser(e: SubmitEvent) {
	e.preventDefault();
	isLoading = true;
	updateSubmitBtn(btnAttrs, isLoading);
	clearFormErrors(frmSignin);

	const isValid = validate();

	if (!isValid) {
		handleFormErrors(validationErrors, frmSignin);
		showNotification("error", "Invalid input!");
		isLoading = false;
		updateSubmitBtn(btnAttrs, isLoading);
		return;
	}

	try {
		const res = await fetch(frmSignin.action, {
			method: frmSignin.method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: inputEmail.value.trim(),
				password: inputPassword.value.trim(),
			}),
		});

		if (!res.ok) {
			const { message, errors }: APIResponse<undefined> =
				await res.json();

			errors && handleFormErrors(errors, frmSignin);
			showNotification("error", message);
		} else {
			const { message, data }: APIResponse<undefined> = await res.json();

			console.log("message", message, "data", data);

			showNotification("success", message);
		}
	} catch (error) {
		console.log("Error during sign-in process:", error);
		if (error instanceof Error) showNotification("error", error.message);
	} finally {
		isLoading = false;
		updateSubmitBtn(btnAttrs, isLoading);
	}
}

function validate(): boolean {
	let isValid = true;

	validationErrors.email = [];
	validationErrors.password = [];

	if (!isRequiredInputFilled(inputEmail)) {
		validationErrors.email.push("Email is required");
		isValid = false;
	}

	if (!isRequiredInputFilled(inputPassword)) {
		validationErrors.password.push("Password is required");
		isValid = false;
	}

	if (!isValidEmail(inputEmail.value)) {
		validationErrors.email.push("Email must be a valid email address");
		isValid = false;
	}

	return isValid;
}
