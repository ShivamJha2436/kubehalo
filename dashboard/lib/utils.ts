export function classNames(...args: Array<string | false | null | undefined>) {
    return args.filter(Boolean).join(" ");
}
