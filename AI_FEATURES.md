# 🤖 AI Features Configuration Guide

## Prerequisites

To use AI-powered features (auto-tagging, document summarization, Q&A), you need an OpenAI API key.

## Getting Your OpenAI API Key

1. **Create an OpenAI Account**
   - Go to [https://platform.openai.com/signup](https://platform.openai.com/signup)
   - Sign up or log in to your account

2. **Generate an API Key**
   - Navigate to [https://platform.openai.com/api-keys](https://platform.openai.com/api-keys)
   - Click "Create new secret key"
   - Give it a name (e.g., "Research Institute System")
   - Copy the key immediately (you won't be able to see it again!)

3. **Add Credits to Your Account**
   - Go to [https://platform.openai.com/account/billing](https://platform.openai.com/account/billing)
   - Add a payment method
   - Add credits (minimum $5 recommended)

## Configuration

### Option 1: Environment Variable (Recommended)

**Windows (PowerShell):**
```powershell
$env:OPENAI_API_KEY="sk-your-api-key-here"
```

**Windows (Command Prompt):**
```cmd
set OPENAI_API_KEY=sk-your-api-key-here
```

**Linux/Mac:**
```bash
export OPENAI_API_KEY=sk-your-api-key-here
```

### Option 2: .env File

1. Copy `.env.example` to `.env`
2. Edit the `.env` file and add your API key:
   ```
   OPENAI_API_KEY=sk-your-api-key-here
   ```

## Available AI Features

### 1. 🏷️ Auto-Tagging
**Location:** Document Upload Page → "🤖 Generate Tags (AI)" button

**How to use:**
1. Fill in document name and description
2. Click "🤖 Generate Tags (AI)" button
3. AI will generate 5 relevant keywords based on the document information
4. Tags will be automatically filled in the Keywords field

**Example:**
- Document: "Machine Learning Research Paper"
- Description: "Analysis of deep learning algorithms for image recognition"
- Generated Tags: `machine learning, deep learning, image recognition, neural networks, computer vision`

### 2. 📝 Auto-Description Generation
**Location:** Document Upload Page → "🤖 Generate Description" button

**How to use:**
1. Fill in document name and select document type
2. Optionally select a file
3. Click "🤖 Generate Description" button
4. AI will generate a professional 2-3 sentence description
5. Description will be automatically filled in the Description field

**Example:**
- Document Name: "Q3 2025 Financial Report"
- Document Type: "Report"
- Generated Description: "This quarterly financial report provides a comprehensive analysis of the organization's financial performance for Q3 2025. It includes detailed revenue breakdowns, expense tracking, and key performance indicators. The report serves as an essential resource for stakeholders to evaluate financial health and strategic planning."

### 3. 📄 Document Summarization (Post-Upload)
**API Method:** `GenerateDocumentSummary(documentID, maxLength)`

**How to use:**
```javascript
const summary = await GenerateDocumentSummary(123, 200)
```

### 4. ❓ Ask Questions About Documents
**API Method:** `AskDocumentQuestion(documentID, question)`

**How to use:**
```javascript
const answer = await AskDocumentQuestion(123, "What are the main findings?")
```

## Cost Estimation

Using GPT-3.5 Turbo:
- **Auto-tagging:** ~$0.001 per document
- **Auto-description:** ~$0.001-0.002 per document
- **Summarization:** ~$0.002 per document
- **Q&A:** ~$0.001-0.003 per question

**Example monthly cost for 1000 documents:**
- Auto-tagging: ~$1
- Auto-description: ~$1.50
- Total: Less than $5/month for typical usage

## Troubleshooting

### "OpenAI API key not configured" Error

**Solution:**
1. Make sure you've set the `OPENAI_API_KEY` environment variable
2. Restart the application after setting the variable
3. Check that your API key starts with `sk-`

### "Insufficient quota" Error

**Solution:**
1. Go to [https://platform.openai.com/account/billing](https://platform.openai.com/account/billing)
2. Add credits to your account
3. Try again

### "Rate limit exceeded" Error

**Solution:**
1. Wait a few seconds and try again
2. If persistent, upgrade your OpenAI plan at [https://platform.openai.com/account/limits](https://platform.openai.com/account/limits)

## Security Best Practices

⚠️ **IMPORTANT:**
- Never commit your `.env` file to Git
- Never share your API key publicly
- Rotate your API key regularly
- Set usage limits in OpenAI dashboard to prevent unexpected charges

## Additional Resources

- [OpenAI API Documentation](https://platform.openai.com/docs)
- [OpenAI Pricing](https://openai.com/pricing)
- [OpenAI Usage Dashboard](https://platform.openai.com/usage)
