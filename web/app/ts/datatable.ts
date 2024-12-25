const datatable = document.getElementById("datatable");
const thead = datatable?.querySelector("thead");
const tbody = datatable?.querySelector("tbody");
const endpoint = datatable?.dataset.url;

type TableHeaders = {
	label: string;
	field: string;
};

const tblHeaders: TableHeaders[] = JSON.parse(
	datatable?.dataset.headers || "[]"
);

const tblData = JSON.parse(datatable?.dataset.data || "[]");

function renderTHead() {
	const fragment = document.createDocumentFragment();
	const row = document.createElement("tr");
	tblHeaders.forEach((data) => {
		const th = document.createElement("th");
		th.textContent = data.label;
		row.appendChild(th);
	});
	fragment.appendChild(row);
}

function renderTBody() {
	const fragment = document.createDocumentFragment();
	const row = document.createElement("tr");
	fragment.appendChild(row);
}

renderTHead();
renderTBody();
