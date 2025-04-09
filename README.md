# Erebrus Beacon Node

The Erebrus Beacon Node acts as a decentralized VPN relay, providing censorship-resistant and secure internet access. It allows users to route their traffic privately through the Erebrus network while helping strengthen the network's resilience.

## Overview

By operating a Beacon Node, you help expand decentralized connectivity while preventing surveillance, censorship, and cyber threats. The setup is lightweight and runs in a Docker container, requiring minimal system resources.

## Prerequisites

- Operating System: `Linux`
- Minimum Hardware Requirements:
  - 2GB RAM
  - 1vCPU
- Network Requirements:
  - Incoming traffic allowed on ports:
    - 51820 (UDP)
    - 9080 (TCP)
    - 9002 (TCP)
- A computer with a reliable internet connection (preferably wired)
- Basic understanding of command line interface (CLI) tools
- Node Wallet (valid mnemonic on the supported chain) for on-chain registration and checkpoint

## Deployment Options

### Cloud Deployment (Recommended)

For optimal performance and uptime, cloud deployment is recommended. Here are estimated monthly costs from different providers:

| Traffic Level | Hetzner | AWS | DigitalOcean | Vultr |
|---------------|---------|-----|--------------|-------|
| Low (100GB)   | $5      | $10 | $12         | $8    |
| Moderate (500GB) | $5   | $50 | $17         | $8    |
| High (2TB)    | $5      | $180| $22         | $8    |

### Home Lab Setup

For home lab deployments:
- Requires a static IP from your ISP
- Ensure unlimited or high-bandwidth ISP plan
- Recommended: Use a dedicated VPS or bare-metal server with static IP

## Requirements

- **License Fee:** No license required - open for community participation
- **Staking Requirement:** Currently $0 USD (subject to change based on demand & supply)
- **Revenue Model:** Node operators earn rewards based on:
  - Bandwidth contribution
  - Uptime
  - Reliability

## Installation

1. **Install Node Software**
   ```bash
   # If running as regular user:
   sudo bash <(curl -s https://raw.githubusercontent.com/NetSepio/beacon/main/install.sh)
   
   # If running as root:
   bash <(curl -s https://raw.githubusercontent.com/NetSepio/beacon/main/install.sh)
   ```

2. **Configure Node Parameters**
   - Installation Directory (default: current working directory)
   - Public IP (must be publicly routable)
   - Chain selection
   - Valid mnemonic

3. **Verification Process**
   - Configuration verification
   - Public IP reachability test
   - Service startup

## Maintenance & Monitoring

- Regular software updates for security patches and performance improvements
- Monitor node health and resource usage
- Maintain stable internet connection

## Security Considerations

- Protect your Node Operator account mnemonic
- Keep node software and OS updated
- Implement best practices for server security
- Consider DDoS protection measures

## Additional Resources

- [Erebrus Documentation](https://docs.netsepio.com/latest/erebrus/nodes/beacon-node)
- [Support Discord](https://discord.gg/netsepio)

## License

This project is open-source and available for community participation.
