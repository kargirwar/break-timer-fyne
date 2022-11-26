import { Utils } from './modules/utils.js';
class Index {
    constructor() {
        document.addEventListener('DOMContentLoaded', () => {
            this.initDom();
            this.initHandlers();
        });
    }

    initDom() {
        this.$addRule = document.getElementById('add-rule');
        this.$saveRules = document.getElementById('save-rules');
        this.$rules = document.getElementById('rules');
        this.ruleTemplate = document.getElementById('rule-template').innerHTML;
    }

    initHandlers() {
        this.$addRule.addEventListener('click', () => {
            const n = Utils.generateNode(this.ruleTemplate, {});
            this.$rules.append(n);
        });

        this.$saveRules.addEventListener('click', () => {
            try {
                this.validate();
            } catch (e) {
                return;
            }
            this.saveRules();
        });
    }

    validate() {
    }

    saveRules() {
        const $rules = [...document.querySelectorAll('.rule-container')];
        let rules = {
            rules: []
        };

        $rules.forEach(($r) => {
            const interval = $r.querySelector('.interval').value;
            const days = this.getDays($r.querySelector('.days').value);
            const from = $r.querySelector('.from').value;
            const to = $r.querySelector('.to').value;

            console.log(`${interval} ${days} ${from} ${to}`);
            rules.rules.push({
                days,
                interval,
                from,
                to,
            });
        });

        console.log(JSON.stringify(rules));
    }

    getDays(type) {
        console.log(type);
        switch (type) {
            case 'weekdays':
                return ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"];

            case 'weekends':
                return ["Saturday", "Sunday"];
        }
    }
}

new Index();
