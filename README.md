# Stasyan

Stasyan is very simple chaos enginering tool. By default it just every 5 min delete random pod in default k8s namespace for a 1 hour. Main idea is in an hour you can check that your app has desired DR ability

## Usage

I simply do smth like 

```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: stasyan      
spec:
  containers:
  - name: stasyan
    image: ondator/stasyan
    env: 
    - name: STASYAN_LIFETIME
      value: "10"
    - name: MY_POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
EOF
```

## Configuration

Now it's just a couple of envs:

`STASYAN_NAMESPACE` for declaring namespace where Stasyan will crush
`STASYAN_LIFETIME` for declaring time (in minuts) during Stasyan will crush

## WTF mean Stasyan?!

It's russian meme about clumsy guy named Stasyan: ![](https://cs14.pikabu.ru/post_img/2022/11/09/5/1667974967196750168.webp)