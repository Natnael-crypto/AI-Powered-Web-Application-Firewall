import React, { useState, useEffect } from "react";

type Condition = {
  ruleType: string;
  ruleMethod: string;
  ruleDefinition: string;
};

type RuleInput = {
  ruleID: string;
  action: string;
  category: string;
  conditions: Condition[];
  applications: string[];
};

const availableApps = [
  "App-One",
  "App-Two",
  "App-Three",
  "Customer-Portal",
  "Internal-API",
];

const validRuleTypes = [
  "REQUEST_HEADERS", "REQUEST_URI", "ARGS", "ARGS_GET", "ARGS_POST",
  "REQUEST_COOKIES", "REQUEST_BODY", "XML", "JSON", "REQUEST_METHOD",
  "REQUEST_PROTOCOL", "REMOTE_ADDR",
];

const validRuleMethods = [
  "regex", "streq", "contains", "ipMatch", "rx", "beginsWith",
  "endsWith", "eq", "pm",
];

const validActions = [
  "deny", "log", "nolog", "pass", "drop", "redirect", "capture",
  "t:none", "t:lowercase", "t:normalizePath", "t:urlDecode",
  "t:compressWhitespace", "severity:2", "severity:3", "status:403",
];

const CreateRuleModal: React.FC = () => {
  const [ruleInput, setRuleInput] = useState<RuleInput>({
    ruleID: "1001",
    action: "deny",
    category: "SQL Injection",
    conditions: [{
      ruleType: "ARGS",
      ruleMethod: "contains",
      ruleDefinition: "SELECT",
    }],
    applications: [],
  });

  const [preview, setPreview] = useState<string>("");

  useEffect(() => {
    generateRule(ruleInput);
  }, [ruleInput]);

  const updateCondition = (index: number, field: keyof Condition, value: string) => {
    const updatedConditions = [...ruleInput.conditions];
    updatedConditions[index][field] = value;
    setRuleInput({ ...ruleInput, conditions: updatedConditions });
  };

  const handleAppAdd = (app: string) => {
    if (!ruleInput.applications.includes(app)) {
      setRuleInput({ ...ruleInput, applications: [...ruleInput.applications, app] });
    }
  };

  const handleAppRemove = (app: string) => {
    setRuleInput({
      ...ruleInput,
      applications: ruleInput.applications.filter(a => a !== app),
    });
  };

  const generateRule = (input: RuleInput) => {
    const { ruleID, action, category, conditions } = input;
    if (conditions.length === 0) return;

    let ruleText = "";
    conditions.forEach((cond, i) => {
      const prefix = i === 0 ? "SecRule" : "    SecRule";
      const chain = i < conditions.length - 1 ? `"chain"` : "";
      const firstLine = i === 0
        ? `"id:${ruleID},phase:2,${action},msg:'${category}'${conditions.length > 1 ? ",chain" : ""}"`
        : chain;
      ruleText += `${prefix} ${cond.ruleType} "@${cond.ruleMethod} ${cond.ruleDefinition}" ${firstLine}\n`;
    });

    setPreview(ruleText.trim());
  };

  const addCondition = () => {
    setRuleInput({
      ...ruleInput,
      conditions: [
        ...ruleInput.conditions,
        { ruleType: "", ruleMethod: "", ruleDefinition: "" },
      ],
    });
  };

  const removeCondition = (index: number) => {
    const updated = [...ruleInput.conditions];
    updated.splice(index, 1);
    setRuleInput({ ...ruleInput, conditions: updated });
  };

  const unselectedApps = availableApps.filter(app => !ruleInput.applications.includes(app));

  return (
    <div className="p-4 space-y-4 max-w-4xl mx-auto">
      <h2 className="text-xl font-bold">WAF Rule Generator</h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <input
          className="border p-2"
          placeholder="Rule ID"
          value={ruleInput.ruleID}
          onChange={(e) => setRuleInput({ ...ruleInput, ruleID: e.target.value })}
        />

        <select
          className="border p-2"
          value={ruleInput.action}
          onChange={(e) => setRuleInput({ ...ruleInput, action: e.target.value })}
        >
          <option value="">Select Action</option>
          {validActions.map(action => (
            <option key={action} value={action}>{action}</option>
          ))}
        </select>

        <input
          className="border p-2"
          placeholder="Category"
          value={ruleInput.category}
          onChange={(e) => setRuleInput({ ...ruleInput, category: e.target.value })}
        />
      </div>

      {/* Application Dropdown */}
      <div>
        <label className="block font-semibold mb-2">Select Application</label>
        <select
          className="border p-2 w-full"
          onChange={(e) => {
            const value = e.target.value;
            if (value) handleAppAdd(value);
            e.target.value = "";
          }}
        >
          <option value="">-- Select Application --</option>
          {unselectedApps.map(app => (
            <option key={app} value={app}>{app}</option>
          ))}
        </select>

        <div className="flex flex-wrap gap-2 mt-3">
          {ruleInput.applications.map(app => (
            <div
              key={app}
              className="bg-green-200 text-green-900 px-3 py-1 rounded-full cursor-pointer transition duration-200 hover:bg-red-200 hover:text-red-900"
              onClick={() => handleAppRemove(app)}
              title="Click to remove"
            >
              {app}
            </div>
          ))}
        </div>
      </div>

      {/* Conditions */}
      {ruleInput.conditions.map((cond, index) => (
        <div key={index} className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center mb-2">
          <select
            className="border p-2"
            value={cond.ruleType}
            onChange={(e) => updateCondition(index, "ruleType", e.target.value)}
          >
            <option value="">Select Rule Type</option>
            {validRuleTypes.map(type => (
              <option key={type} value={type}>{type}</option>
            ))}
          </select>

          <select
            className="border p-2"
            value={cond.ruleMethod}
            onChange={(e) => updateCondition(index, "ruleMethod", e.target.value)}
          >
            <option value="">Select Method</option>
            {validRuleMethods.map(method => (
              <option key={method} value={method}>{method}</option>
            ))}
          </select>

          <div className="flex gap-2">
            <input
              className="border p-2 flex-1"
              placeholder="Definition"
              value={cond.ruleDefinition}
              onChange={(e) => updateCondition(index, "ruleDefinition", e.target.value)}
            />
            {index > 0 && (
              <button
                className="bg-red-500 text-white px-2 py-1 rounded"
                onClick={() => removeCondition(index)}
              >
                ✕
              </button>
            )}
          </div>
        </div>
      ))}

      <button
        className="bg-blue-500 text-white px-4 py-2 rounded"
        onClick={addCondition}
      >
        ➕ Add Condition
      </button>

      {/* Preview */}
      <pre className="bg-gray-100 p-4 rounded border whitespace-pre-wrap mt-4">
        {preview || "// Rule preview will appear here..."}
        {"\n\n"}// Applications: {ruleInput.applications.join(", ") || "None"}
      </pre>
    </div>
  );
};

export default CreateRuleModal;
