from kafka import KafkaConsumer
consumer=KafkaConsumer("Gubarticle_Exchange",group_id="chenhuachao20180226",
                       bootstrap_servers="172.16.56.137:9092,172.16.56.138:9092,172.16.56.139:9092",
                       auto_offset_reset="earliest",
                       enable_auto_commit=True)
print("xxxxxxxx")
for message in consumer:
    print(">>>>",message)