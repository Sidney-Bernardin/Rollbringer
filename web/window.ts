import { Alpine } from "alpinejs";
import { HtmxExtension } from "htmx.org";

declare global {
    interface Window {
        pageData: any;
        alpine: Alpine;
        htmx: HtmxExtension;
    }
}

