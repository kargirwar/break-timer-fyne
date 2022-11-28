const DISABLED = [
]

class Logger {
    static log(tag, str) {
        if (DISABLED.includes(tag)) {
            return;
        }

        Logger.print(tag, str);
    }

    static print(tag, str) {
        let [month, date, year]    = new Date().toLocaleDateString("en-US").split("/")
        let [hour, minute, second] = new Date().toLocaleTimeString("en-US").split(/:| /)

        let o = `${date}-${month}-${year} ${hour}:${minute}:${second}:::${tag}: ${str}`;
        console.log(o)
    }
}

export { Logger }
