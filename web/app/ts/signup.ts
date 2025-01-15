import { clearFormErrors, showFormError, updateSubmitBtn } from "./form";
import { showNotification } from "./notification";
import {
	isRequiredInputFilled,
	isValidEmail,
	resetValidationErrors,
} from "./validation";

const frm = document.getElementById("frmSignup") as HTMLFormElement;
const inputEmail = frm.querySelector("#email") as HTMLInputElement;
const inputPassword = frm.querySelector("#password") as HTMLInputElement;
const inputRetypePass = frm.querySelector(
	"#password_confirmation"
) as HTMLInputElement;
const btnSignup = frm.querySelector("#btnSignup") as HTMLButtonElement;

const btnAttrs = {
	btn: btnSignup,
	text: "Sign Up",
	loadingText: "Signing up...",
};

const validationErrors: ValidationErrorMap = {
	email: [],
	password: [],
	passwordConfirm: [],
};

let isLoading = false;

function validate(): boolean {
	const email = inputEmail.value;
	const password = inputPassword.value;
	const passwordConfirm = inputRetypePass.value;

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

async function signUpUser(e: SubmitEvent) {
	e.preventDefault();
	isLoading = true;
	updateSubmitBtn(btnAttrs, isLoading);
	clearFormErrors(frm);

	const isValid = validate();

	if (!isValid) {
		if (validationErrors.email.length > 0)
			showFormError(inputEmail, validationErrors.email);
		if (validationErrors.password.length > 0)
			showFormError(inputPassword, validationErrors.password);

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
				password_confirmation: inputRetypePass.value,
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

function comparePasswords(): void {
	const helpText = inputPassword.nextElementSibling as HTMLElement;

	if (
		inputRetypePass.value &&
		inputPassword.value &&
		inputRetypePass.value != inputPassword.value
	) {
		inputPassword.classList.add("error");
		helpText.textContent = "Passwords do not match";
		helpText.style.display = "block";
	} else {
		inputPassword.classList.remove("error");
		helpText.textContent = "";
		helpText.style.display = "none";
	}
}

frm.addEventListener("change", function (event) {
	const target = event.target as HTMLInputElement;
	if (
		target.matches("#password") ||
		target.matches("#password_confirmation")
	) {
		comparePasswords();
	} else if (target.matches("#email")) {
		const helpTxt = target.nextElementSibling as HTMLElement;

		if (helpTxt) {
			helpTxt.classList.remove("error");
			helpTxt.textContent = "";
			helpTxt.style.display = "none";
		}

		target.classList.remove("error");
		if (!isValidEmail(target.value)) {
			showFormError(target, ["Email must be a valid email address"]);
		}
	}
});

frm.addEventListener("submit", signUpUser);

console.log("Waiting for user to sign up...");
