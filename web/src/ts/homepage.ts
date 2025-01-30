import "./main.ts";

const superSectionList = document.querySelectorAll(
    "section[class*=super-sect]",
);
for (const superSection of superSectionList) {
    superSection.setAttribute("hx-trigger", "intersect");
    superSection.addEventListener("htmx:trigger", function () {
        superSection.classList.add("appeared");
    });
}
