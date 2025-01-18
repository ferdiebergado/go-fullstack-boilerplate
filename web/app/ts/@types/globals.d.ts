type APIResponse<T> = {
	message: string;
	errors?: ValidationErrorMap;
	data?: T;
};

type Errors = string[];

type ValidationError = {
	field: string;
	errors: Errors;
};

type ValidationErrors = {
	errors: ValidationError;
};

type NotificationType = "success" | "error";

type ButtonAttrs = {
	btn: HTMLButtonElement;
	text: string;
	loadingText: string;
};

type ValidationErrorMap = {
	[key: string]: Errors;
};
