document.addEventListener("DOMContentLoaded", function () {
    const codeBlocks = document.querySelectorAll(".markdown-content pre");

    codeBlocks.forEach(function (pre) {
        const code = pre.querySelector("code");

        if (!code) {
            return;
        }

        const wrapper = document.createElement("div");
        wrapper.className = "code-block";

        pre.parentNode.insertBefore(wrapper, pre);
        wrapper.appendChild(pre);

        const button = document.createElement("button");

        button.type = "button";
        button.className = "btn btn-sm btn-light code-copy";
        button.textContent = "Copy";

        button.addEventListener("click", async function () {
            await navigator.clipboard.writeText(code.innerText);

            button.textContent = "Copied";

            setTimeout(function () {
                button.textContent = "Copy";
            }, 1500);
        });

        wrapper.appendChild(button);
    });
});