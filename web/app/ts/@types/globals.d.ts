type APIResponse<T> = {
	message: string;
	errors: ValidationError;
	data: T;
};

type ValidationError = {
	field: string;
	errors: string[];
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

type Errors = string[];

type ValidationErrorMap = {
	[key: string]: Errors;
};
