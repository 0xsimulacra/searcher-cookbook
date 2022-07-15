# searcher cookbook
The cookbook is a collection of usefully scripts, functions, smart contracts, helpers and tips for golang evm searchers.  
I will extend the list from time to time with some alpha or just some hints to help you in your searcher journey.  
First I want to apologize for my english, I'm not a native speaker ;).  
Since my journey as a searcher started I had lots of ups and downs. 
All in all you need to be a pretty competitive guy to stay up with all the excellent searchers out there.  
This collection will expect that you have some basic knowledge as I'm not going to explain everything in detail. Still I hope it will help you to get faster into searching.

# Contents:
- [eth_call with state override](eth_call-state-override/README.md)

# Bonus
I will release packages that I use in most of my bots and will give some insights how I use them
- [discogo](https://github.com/tripolious/discogo) - used for sending special logs to my discord channel and update different configurations like miner fee % without restarting the bot or check flashbots_getBundleStats directly from discord.    


### To be continued
At the moment I'm thinking of releasing information based on the following topics.  
I will release them in a way like I did it in the first example.  
Please let me know if you want some other topics that will help beginners or advanced searchers.
```
    - server configuration
    - use discord to log and update your bot
    - utilize eth logs for faster updates to your bot    
    - utilize eth_callBundle, eth_sendBundle and flashbots_getBundleStats in your logs  
    - prepare token lists
    - prepare pair lists
    - uni2 / uni3 onchain calculations
    - uni2 / uni3 arb contracts
    - uni2 offchain calculations
    - uni3 offchain calculations
    - opensea-seaport api
    - looksrare api
    - x2y2 api
    - nft arb contract
    - keep your key's away from your server - especially helpful if you work in a team and more than 1 person can access the server
    - ....
```

### Why do you share all this information?
I know that most of the searchers keep their findings and learnings as a secrete (I also did it :D), still I think we need more devs in the market to push crypto.    
Especially in this market situation we need more devs to join crypto and as a searcher you will get in touch and need to learn nearly everything about crypto.  
Flashbots did an amazing job, and so I want to contribute back.  
Still I will not share everything completely in detail, so you will still need to learn the things and get familiar with the different topics in detail ;)

- Smart Contracts security (make them secure, or you will suffer)
- Smart Contract optimization (don't fear assembly or packed data)
- Smart Contract testing (you will need to test against different states, fork and so on)
- Learn to understand different protocols and how to read their code
- Use of the go-ethereum package (signing, nonce handling, ...)
- Get in touch with different evm-chains to check if your alpha will work there too
- Understand nodes (geth, erigon, ...) and keep your server secured
- get a deep dive into your chosen language (go, node, rust, ...) to write your bot

What can motivate more than earning money while you compete in a wild pvp world ;)?