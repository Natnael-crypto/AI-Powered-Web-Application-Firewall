import axios from "axios";
import React, { useEffect, useState } from "react";
import { AppOption, Condition, RuleInput, validActions, validRuleMethods, validRuleTypes } from "../lib/types";
import { useGetApplications } from "../hooks/api/useApplication";
import { createRule } from "../services/rulesApi";

type CreateRuleModalProps = {
  isOpen: boolean;
  onClose: () => void;
};

const CreateRuleModal: React.FC<CreateRuleModalProps> = ({ isOpen, onClose }) => {
  const [ruleInput, setRuleInput] = useState<RuleInput>({
    ruleID: "1001",
    action: "deny",
    category: "Message",
    conditions: [{
      ruleType: "ARGS",
      ruleMethod: "regex",
      ruleDefinition: "value like *select*",
    }],
    applications: [],
  });
  const [preview, setPreview] = useState<string>("");
  const [availableApps, setAvailableApps] = useState<AppOption[]>([]);
  const {data:applications} = useGetApplications()

 useEffect(() => {
  if (applications) {
    const apps = applications.map(app => ({
      application_id: app.application_id,
      application_name: app.application_name,
    }));
    setAvailableApps(apps);
  }
}, [applications]);


  useEffect(() => {
    generateRule(ruleInput);
  }, [ruleInput]);

  const updateCondition = (index: number, field: keyof Condition, value: string) => {
    const updatedConditions = [...ruleInput.conditions];
    updatedConditions[index][field] = value;
    setRuleInput({ ...ruleInput, conditions: updatedConditions });
  };

  const handleAppAdd = (appId: string) => {
    if (!ruleInput.applications.includes(appId)) {
      setRuleInput({ ...ruleInput, applications: [...ruleInput.applications, appId] });
    }
  };

  const handleAppRemove = (appId: string) => {
    setRuleInput({
      ...ruleInput,
      applications: ruleInput.applications.filter(a => a !== appId),
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
      conditions: [...ruleInput.conditions, { ruleType: "", ruleMethod: "", ruleDefinition: "" }],
    });
  };

  const saveRule = async () => {
    if (ruleInput.applications.length === 0) {
      alert("Please select at least one application.");
      return;
    }

    const payloadTemplate = {
      action: ruleInput.action,
      category: ruleInput.category,
      is_active: true,
      conditions: ruleInput.conditions.map((cond) => ({
        rule_type: cond.ruleType,
        rule_method: cond.ruleMethod,
        rule_definition: cond.ruleDefinition,
      })),
      application_ids: ruleInput.applications,
    };

    console.log(payloadTemplate)

    try {
      createRule(payloadTemplate)
      alert("Rule saved successfully for all selected applications!");
      onClose();
    } catch (error) {
      console.error("Error saving rule:", error);
      alert("An error occurred while saving the rule.");
    }
  };

  const unselectedApps = availableApps.filter(app => !ruleInput.applications.includes(app.application_id));

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white rounded-lg shadow-lg p-6 max-w-4xl w-full max-h-[90vh] overflow-y-auto relative">
        <button onClick={onClose} className="absolute top-3 right-3 text-gray-600 hover:text-red-600 text-xl">
          &times;
        </button>

        <h2 className="text-2xl font-bold mb-4">Create WAF Rule</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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

        <div className="mt-4">
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
              <option key={app.application_id} value={app.application_id}>
                {app.application_name}
              </option>
            ))}
          </select>

          <div className="flex flex-wrap gap-2 mt-2">
            {ruleInput.applications.map(appId => {
              const app = availableApps.find(a => a.application_id === appId);
              return (
                <span
                  key={appId}
                  className="bg-green-200 text-green-900 px-3 py-1 rounded-full cursor-pointer hover:bg-red-200 hover:text-red-900"
                  onClick={() => handleAppRemove(appId)}
                >
                  {app?.application_name || appId}
                </span>
              );
            })}
          </div>
        </div>

        {/* Conditions */}
        {ruleInput.conditions.map((cond, index) => (
          <div key={index} className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center my-2">
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
                  onClick={() => {
                    const updated = [...ruleInput.conditions];
                    updated.splice(index, 1);
                    setRuleInput({ ...ruleInput, conditions: updated });
                  }}
                >
                  ✕
                </button>
              )}
            </div>
          </div>
        ))}

        <div className="flex justify-between items-center mt-4">
          <button
            className="bg-blue-500 text-white px-4 py-2 rounded"
            onClick={addCondition}
          >
          Add Condition
          </button>

          <button
            className="bg-green-500 text-white px-4 py-2 rounded"
            onClick={saveRule}
          >
          Save Rule
          </button>
        </div>

        <pre className="bg-gray-100 p-4 rounded border whitespace-pre-wrap mt-6">
          {preview || "// Rule preview will appear here..."}
        </pre>
      </div>
    </div>
  );
};

export default CreateRuleModal;
