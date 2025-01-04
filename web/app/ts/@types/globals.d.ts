type APIResponse<T> = {
	success: string;
	message: string;
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
