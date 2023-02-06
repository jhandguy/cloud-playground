import grpc from 'k6/net/grpc';
import {check, sleep} from 'k6';
import {randomString, uuidv4} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    stages: [
        {target: 10, duration: '20s'},
        {target: 10, duration: '20s'},
        {target: 0, duration: '20s'},
    ],
    thresholds: {
        'checks': ['rate>0.9'],
        'grpc_req_duration{method:CreateItem}': ['p(95)<10000'],
        'grpc_req_duration{method:GetItem}': ['p(95)<10000'],
        'grpc_req_duration{method:DeleteItem}': ['p(95)<10000'],
    },
};

const client = new grpc.Client();
client.load(['../proto/item'], 'item.pg.proto');

export default function () {
    const params = {
        metadata: {
            'x-token': __ENV.DYNAMO_TOKEN,
        },
    };
    const id = uuidv4();
    const content = randomString(50);
    const data = {
        item: {
            id: id,
            content: content,
        }
    };

    client.connect(__ENV.DYNAMO_URL, {
        plaintext: true,
    });

    check(client.invoke('item.ItemService/CreateItem', data, params), {
        'create status is OK': (r) => r.status === grpc.StatusOK,
        'create message is valid': (r) => r.message.item.id === id && r.message.item.content === content,
    });

    sleep(0.5);

    check(client.invoke('item.ItemService/GetItem', {id: id}, params), {
        'get status is OK': (r) => r.status === grpc.StatusOK,
        'get message is valid': (r) => r.message.item.id === id && r.message.item.content === content,
    });

    sleep(0.5);

    check(client.invoke('item.ItemService/DeleteItem', {id: id}, params), {
        'delete status is OK': (r) => r.status === grpc.StatusOK,
    });

    sleep(1);
}
