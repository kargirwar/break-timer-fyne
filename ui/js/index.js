import { Utils } from './modules/utils.js';
import { Logger } from './modules/logger.js';

const RULES = 'rules';
const DEFAULT_INTERVAL = 30;
const DEFAULT_FROM = 9;
const DEFAULT_TO = 17;
const TAG = "index";

class Index {
    constructor() {
        document.addEventListener('DOMContentLoaded', () => {
            this.serial = 1;
            this.initDom();
            this.initHandlers();
            this.renderRules();
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
            this.renderRule();
        });

        this.$saveRules.addEventListener('click', () => {
            try {
                this.validate();
            } catch (e) {
                return;
            }
            this.saveRules();
        });

        this.$rules.addEventListener('click', (e) => {
            if (!e.target.classList.contains('del-rule')) {
                return;
            }

            const id = e.target.id.split("-")[1];
            Logger.log(TAG, `Deleting ${id}`);
            document.getElementById(`${id}`).remove();
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
            const days = this.getDays($r.querySelector('.days-container'));
            const from = $r.querySelector('.from').value;
            const to = $r.querySelector('.to').value;

            Logger.log(TAG, `${interval} ${days} ${from} ${to}`);
            rules.rules.push({
                days,
                interval,
                from,
                to,
            });
        });

        Utils.saveToLocalStorage(RULES, JSON.stringify(rules)); 
    }

    renderRules() {
        const str = Utils.getFromLocalStorage(RULES);
        if (str) {
            const rules = JSON.parse(str);
            rules.rules.forEach((r) => {
                this.renderRule(r);
            });
        }
    }

    renderRule(r = {}) {
        Logger.log(TAG, `serial: ${this.serial} r: ${JSON.stringify(r)}`);
        const n = Utils.generateNode(this.ruleTemplate, {i: this.serial++});
        n.querySelector('.interval').value = r.interval ?? DEFAULT_INTERVAL;

        if (r.days) {
            [...n.querySelectorAll("[name='days']")].map(($d) => {
                if (r.days.includes($d.value)) {
                    $d.checked = true;
                } else {
                    $d.checked = false;
                }
            });
        }

        n.querySelector('.from').value = r.from ?? DEFAULT_FROM
        n.querySelector('.to').value = r.to ?? DEFAULT_TO

        this.$rules.append(n);
    }

    getDays($container) {
        const $cbs = [...$container.querySelectorAll("[name='days']")];
        let days = [];
        $cbs.forEach(($cb) => {
            if ($cb.checked) {
                days.push($cb.value);
            }
        });

        return days;
    }
}

new Index();
