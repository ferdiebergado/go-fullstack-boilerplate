import {
	clearFormErrors,
	handleFormErrors,
	showFormError,
	toggleError,
	updateSubmitBtn,
} from "./form";
import { showNotification } from "./notification";
import { isRequiredInputFilled, isValidEmail } from "./validation";

const frmSignup = document.getElementById("frmSignup") as HTMLFormElement;
const inputEmail = frmSignup.querySelector("#email") as HTMLInputElement;
const inputPassword = frmSignup.querySelector("#password") as HTMLInputElement;
const inputRetypePass = frmSignup.querySelector(
	"#password_confirmation"
) as HTMLInputElement;
const btnSignup = frmSignup.querySelector("#btnSignup") as HTMLButtonElement;

const btnAttrs = {
	btn: btnSignup,
	text: "Sign Up",
	loadingText: "Signing up...",
};

type ValidationErrors = {
	email: Errors;
	password: Errors;
	passwordConfirm: Errors;
};

const validationErrors: ValidationErrors = {
	email: [],
	password: [],
	passwordConfirm: [],
};

let isLoading = false;

frmSignup.addEventListener("change", handleInputChange);
frmSignup.addEventListener("submit", signUpUser);

function handleInputChange(event: Event) {
	const target = event.target as HTMLInputElement;
	if (
		target.matches("#password") ||
		target.matches("#password_confirmation")
	) {
		comparePasswords();
	} else if (target.matches("#email")) {
		toggleError(event.target as HTMLInputElement, "");

		if (!isValidEmail(target.value)) {
			showFormError(target, ["Email must be a valid email address"]);
		}
	}
}

async function signUpUser(e: SubmitEvent) {
	e.preventDefault();
	isLoading = true;
	updateSubmitBtn(btnAttrs, isLoading);
	clearFormErrors(frmSignup);

	const isValid = validate();

	if (!isValid) {
		handleFormErrors(validationErrors, frmSignup);
		showNotification("error", "Invalid input!");
		isLoading = false;
		updateSubmitBtn(btnAttrs, isLoading);
		return;
	}

	try {
		const res = await fetch(frmSignup.action, {
			method: frmSignup.method,
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: inputEmail.value,
				password: inputPassword.value,
				password_confirmation: inputRetypePass.value,
			}),
		});

		if (!res.ok) {
			const { message, errors }: APIResponse<undefined> =
				await res.json();

			errors && handleFormErrors(errors, frmSignup);
			showNotification("error", message);
		} else {
			const { message, data }: APIResponse<undefined> = await res.json();

			console.log("message", message, "data", data);

			showNotification("success", message);
		}
	} catch (error) {
		console.log("Error during sign-up process:", error);
		if (error instanceof Error) showNotification("error", error.message);
	} finally {
		isLoading = false;
		updateSubmitBtn(btnAttrs, isLoading);
	}
}

function validate(): boolean {
	const email = inputEmail.value;
	const password = inputPassword.value;
	const passwordConfirm = inputRetypePass.value;

	validationErrors.email = [];
	validationErrors.password = [];
	validationErrors.passwordConfirm = [];

	let isValid = true;

	if (!isRequiredInputFilled(inputEmail)) {
		validationErrors.email.push("Email is required");
		isValid = false;
	}

	if (!isRequiredInputFilled(inputPassword)) {
		validationErrors.password.push("Password is required");
		isValid = false;
	}

	if (!isRequiredInputFilled(inputRetypePass)) {
		validationErrors.passwordConfirm.push(
			"Password confirmation is required"
		);
		isValid = false;
	}

	if (!isValidEmail(email)) {
		validationErrors.email.push("Email must be a valid email address");
		isValid = false;
	}

	if (password && passwordConfirm && password != passwordConfirm) {
		validationErrors.password.push("Passwords do not match");
		isValid = false;
	}

	return isValid;
}

function comparePasswords(): void {
	if (
		inputRetypePass.value &&
		inputPassword.value &&
		inputRetypePass.value !== inputPassword.value
	) {
		toggleError(inputPassword, "Passwords do not match");
	} else {
		toggleError(inputPassword, "");
	}
}
