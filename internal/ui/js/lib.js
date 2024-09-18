(function() {
    const TOAST_TTL = 4000

    let container = document.querySelector(".toaster-container")
    if (!container) {
        console.warn("didn't find toaster container")
        return
    }

    function createToast(msg, title) {
        let d = document.createElement("li")
        d.className = "mb-4 group pointer-events-auto relative flex w-full items-center justify-between space-x-2 overflow-hidden rounded-md p-4 pr-6 shadow-lg transition-all data-[swipe=cancel]:translate-x-0 data-[swipe=end]:translate-x-[var(--radix-toast-swipe-end-x)] data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)] data-[swipe=move]:transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[swipe=end]:animate-out data-[state=closed]:fade-out-80 data-[state=closed]:slide-out-to-right-full data-[state=open]:slide-in-from-top-full data-[state=open]:sm:slide-in-from-bottom-full border bg-background text-foreground"
        d.setAttribute("style", "user-select: none; touch-action: none;")
        d.setAttribute("data-state", "closed")
        d.setAttribute("data-swipe-direction", "right")
        d.addEventListener("click", function() {
            removeToast(d)
        })

        let content = document.createElement("div")
        content.className = "grid gap-1"
        d.appendChild(content)

        if (title) {
            let titleEl = document.createElement("div")
            titleEl.className = "text-sm font-semibold [&amp;+div]:text-xs"
            titleEl.innerHTML = title
            content.appendChild(titleEl)
        }

        let msgEl = document.createElement("div")
        msgEl.className = "text-sm opacity-90"
        msgEl.innerHTML = msg
        content.appendChild(msgEl)
        
        container.appendChild(d)

        setTimeout(function() {
            d.setAttribute("data-state", "open")
        }, 100)

        setTimeout(function() {
            removeToast(d)
        }, TOAST_TTL)
    }

    function removeToast(t) {
        //TODO: prevent double removal
        t.setAttribute("data-state", "closed")
        setTimeout(function() {
            container.removeChild(t)
        }, 500)
    }

    document.body.addEventListener("toast", function(evt) {
        if (evt.detail.message) {
            createToast(evt.detail.message, evt.detail.title)
        }
    })

    /**
     * @param {Element} elt 
     */
    function registerClipboardHandlers(elt) {
        elt.querySelectorAll("[data-clipboard]").forEach((e) => {
            const content = e.getAttribute("data-clipboard")
            e.addEventListener("click", () => {
                navigator.clipboard.writeText(content)
            })
        })
    }
    document.addEventListener("htmx:load", (e) => registerClipboardHandlers(e.detail.elt))

    /**
     * 
     * @param {Element} elt 
     */
    function registerCloseHandlers(elt) {
        elt.querySelectorAll("[data-remove]").forEach((e) => {
            const target = e.getAttribute("data-remove")
            e.addEventListener("click", () => {
                document.getElementById(target).remove()
            })
        })
    }
    document.addEventListener("htmx:load", (e) => registerCloseHandlers(e.detail.elt))

})()