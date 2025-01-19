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

type ValidationErrorMap = {
	[key: string]: Errors;
};

type NotificationType = "success" | "error";

type ButtonAttrs = {
	btn: HTMLButtonElement;
	text: string;
	loadingText: string;
};

type RedirectData = {
	redirectUrl: string;
};
