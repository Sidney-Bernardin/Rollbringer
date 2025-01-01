import { Alpine } from "alpinejs";
import { HtmxExtension } from "htmx.org";

declare global {
    interface Window {
        utils: any;
        alpine: Alpine;
        htmx: HtmxExtension;
    }
}

