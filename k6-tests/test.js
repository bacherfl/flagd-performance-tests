import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

export const options = {
    stages: [
        { target: 10, duration: '3m' },
        { target: 50, duration: '3m' },
        { target: 100, duration: '3m' },
        { target: 200, duration: '3m' },
        { target: 300, duration: '3m' },
        { target: 500, duration: '3m' },
        { target: 1000, duration: '3m' },
        { target: 2000, duration: '3m' },
        { target: 3000, duration: '3m' },
        { target: 5000, duration: '3m' },
        { target: 10000, duration: '3m' },
        { target: 50000, duration: '3m' },
    ],
};
export default function () {
    const url = 'http://flagd.flagd-performance-test:80/schema.v1.Service/ResolveString';

    var randomNumber = Math.floor(Math.random() * 5000);

    const payload = {
        flagKey: 'color-' + randomNumber,
        context: {
            "version": "1.0.0"
        }
    };
    const headers = {
        'Content-Type': 'application/json'
    };

    const response = http.post(url, JSON.stringify(payload), { headers });

    check(response, {
        'http response status code is 200': response.status === 200,
    });

    sleep(10)
}