import grpc from 'k6/net/grpc';
import {check} from 'k6';
import {randomString, uuidv4} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    scenarios: {
        load: {
            executor: 'ramping-arrival-rate',
            startRate: 1,
            timeUnit: '1s',
            preAllocatedVUs: 10,
            stages: [
                {target: 10, duration: '40s'},
                {target: 0, duration: '20s'},
            ],
        },
    },
    thresholds: {
        'checks': ['rate>0.9'],
        'grpc_req_duration{method:CreateObject}': ['p(95)<10000'],
        'grpc_req_duration{method:GetObject}': ['p(95)<10000'],
        'grpc_req_duration{method:DeleteObject}': ['p(95)<10000'],
    },
};

const client = new grpc.Client();
client.load(['../proto/object'], 'object.pg.proto');

export default function () {
    const params = {
        metadata: {
            'x-token': __ENV.S3_TOKEN,
        },
    };
    const id = uuidv4();
    const content = randomString(50);
    const data = {
        object: {
            id: id,
            content: content,
        }
    };

    client.connect(__ENV.S3_URL, {
        plaintext: true,
    });

    check(client.invoke('object.ObjectService/CreateObject', data, params), {
        'create status is OK': (r) => r.status === grpc.StatusOK,
        'create message is valid': (r) => r.message.object.id === id && r.message.object.content === content,
    });

    check(client.invoke('object.ObjectService/GetObject', {id: id}, params), {
        'get status is OK': (r) => r.status === grpc.StatusOK,
        'get message is valid': (r) => r.message.object.id === id && r.message.object.content === content,
    });

    check(client.invoke('object.ObjectService/DeleteObject', {id: id}, params), {
        'delete status is OK': (r) => r.status === grpc.StatusOK,
    });
}
