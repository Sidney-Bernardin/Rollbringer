import { Alpine } from "alpinejs";
import { HtmxExtension } from "htmx.org";

declare global {
  interface Window {
    alpine: Alpine;
    htmx: HtmxExtension;
  }
}
