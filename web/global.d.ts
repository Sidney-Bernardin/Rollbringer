export { };

declare global {
    interface Window {
        googleLoginCallback: (data: any) => void;
    }
}
