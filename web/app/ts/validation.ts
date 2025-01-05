export function isValidEmail(email: string) {
	const emailPattern =
		/^(?i:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$/; // Simple email regex
	return emailPattern.test(email);
}

export function isRequiredInputFilled(input: HTMLInputElement) {
	return input.value.trim() !== ""; // Check for non-empty and non-whitespace value
}
