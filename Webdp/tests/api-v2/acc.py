URL_ACCURACY_F = lambda engine: f"{BASE}/queries/accuracy?engine={engine}"

query_build = lambda q: {
    "dataset" : 1,
    "budget"  : {"epsilon" : 0.05},
    "query"   : q
}
query = [
    {
        "filter" : ["job_sector == \"Finance\""]
    },
    {
        "count" : {
            "column" : "job_sector",
            "mech" : "Laplace"
        }
    }
]
query2 = [
    {
        "filter" : ["job_sector == \"Finance\""]
    },
    {
        "sum" : {
            "column" : "salary_SEK",
            "mech" : "Laplace"
        }
    }
]


# lets get an estimate of how accurate the mean value is

confidence = math.sqrt(0.95)

acc_q = lambda engine, q: requests.post(
    url=URL_ACCURACY_F(engine=engine),
    json=q,
    headers=auth(analyst_token)
)

query_a  = query_build(query)
query_a["confidence"] = confidence

query2_a = query_build(query2)
query2_a["confidence"] = confidence

count_acc = acc_q("opendp", query_a)
sum_acc   = acc_q("opendp", query2_a)

count_acc_res = float(count_acc.json()[0])
sum_acc_res   = float(sum_acc.json()[0])

print(f"Accuracy of count {count_acc_res} at {confidence} level of confidence")
print(f"Accuracy of sum {sum_acc_res} at {confidence} level of confidence")