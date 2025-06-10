#!/bin/bash

# AI vs Template Mode Comparison Test
# Run this script with your OpenAI API key to see the differences

if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ Please set your OpenAI API key:"
    echo "export OPENAI_API_KEY=your-actual-key"
    echo "Then run: ./ai-comparison-test.sh"
    exit 1
fi

echo "🧪 Crossplane AI: Template Mode vs Real AI Comparison"
echo "======================================================"

# Test scenarios to demonstrate AI advantages
SCENARIOS=(
    "I have a database that keeps failing to connect and pods restart every few minutes. What debugging steps should I take?"
    "My AWS RDS instance costs are getting too high. Can you analyze my setup and suggest optimizations?"
    "I'm getting intermittent 502 errors from my application. The load balancer seems fine but something is wrong with the backend."
    "I want to migrate from AWS to GCP. What Crossplane resources do I need to modify and in what order?"
    "My Kubernetes cluster is running out of capacity but I can't figure out which resources are consuming the most."
)

for i in "${!SCENARIOS[@]}"; do
    scenario="${SCENARIOS[$i]}"
    echo -e "\n🔍 Test Scenario $((i+1)): Complex Infrastructure Question"
    echo "Question: $scenario"
    echo ""
    
    echo "📝 TEMPLATE MODE Response:"
    echo "----------------------------------------"
    ./crossplane-ai ask "$scenario" | head -10
    echo ""
    
    echo "🤖 AI MODE Response:"
    echo "----------------------------------------"
    OPENAI_API_KEY=$OPENAI_API_KEY ./crossplane-ai --config ./config.yaml ask "$scenario" | head -15
    echo ""
    echo "========================================================"
done

echo -e "\n🎯 Key Differences You Should Notice:"
echo "📝 Template Mode: Generic, pattern-based responses"
echo "🤖 AI Mode: Specific, contextual, actionable advice"
echo ""
echo "💡 The AI mode provides:"
echo "   • Detailed troubleshooting steps"
echo "   • Specific kubectl commands" 
echo "   • Root cause analysis"
echo "   • Tailored recommendations"
echo "   • Context-aware solutions"
