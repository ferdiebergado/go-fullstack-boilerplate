const notification = document.querySelector(".notification") as HTMLElement;

export function showNotification(notifType: NotificationType, msg: string) {
	switch (notifType) {
		case "success":
			if (notification.classList.contains("error")) {
				notification.classList.replace("error", notifType);
			} else {
				notification.classList.add(notifType);
			}
			notification.textContent = "SUCCESS: ";
			break;

		case "error":
			if (notification.classList.contains("success")) {
				notification.classList.replace("success", notifType);
			} else {
				notification.classList.add(notifType);
			}
			notification.textContent = "ERROR: ";
			break;
	}

	notification.textContent += msg;
	notification.style.display = "block";
}
