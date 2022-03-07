import http from 'k6/http';
import {check, sleep} from 'k6';
import {randomString, uuidv4} from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

export const options = {
    stages: [
        {target: 40, duration: '1m'},
        {target: 40, duration: '1m'},
        {target: 0, duration: '1m'},
    ],
    thresholds: {
        'checks': ['rate>0.9'],
        'http_req_duration{method:POST}': ['p(95)<10000'],
        'http_req_duration{method:GET}': ['p(95)<10000'],
        'http_req_duration{method:DELETE}': ['p(95)<10000'],
    },
};

export default function () {
    const url = `http://${__ENV.GATEWAY_URL}`;
    const params = {
        headers: {
            'Authorization': __ENV.GATEWAY_TOKEN,
            'Host': __ENV.GATEWAY_HOST,
            'Content-Type': 'application/json',
        },
    };
    const id = uuidv4()
    const content = randomString(50)
    const data = JSON.stringify({
        id: id,
        content: content,
    });

    check(http.post(`${url}/message`, data, params), {
        'post response status is 200': (r) => r.status === 200,
        'post response body is valid': (r) => r.json().id === id && r.json().content === content,
    });

    sleep(0.5)

    check(http.get(`${url}/message/${id}`, params), {
        'get response status is 200': (r) => r.status === 200,
        'get response body is valid': (r) => r.json().id === id && r.json().content === content,
    });

    sleep(0.5)

    check(http.del(`${url}/message/${id}`, null, params), {
        'delete response status is 200': (r) => r.status === 200,
    });

    sleep(1)
}
