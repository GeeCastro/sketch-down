# Boat Schematics App: Component Library Plan

This document outlines the first set of specific components to include in the boat schematics app's component library. It includes basic electrical specifications, links to diagram-friendly images, and a mapping of these real-world specs to our JSON schema.

## 1. Inverters

### Sterling Combi Pro S
* **Description**: Pure sine wave combined inverter-charger.
* **Voltage**: 12V or 24V DC input, 110V/220V AC output.
* **Max Amps (Charging)**: 80A (12V/2500W), 100A (12V/3500W).
* **Capacity/Power**: 2500W (2100W continuous) or 3500W (3200W continuous).
* **Key Ports/Connections**: DC battery terminals, AC input, AC output, remote control port, earth-neutral link.
* **Diagram-Friendly Image**: [Sterling Combi Pro S](https://sterling-power.com/cdn/shop/products/PCS122500.jpg) *(Example)*

### Victron MultiPlus
* **Description**: Combined true sine wave inverter and adaptive battery charger.
* **Voltage**: 12V, 24V, or 48V DC input; 120V or 230V AC output.
* **Max Amps (Charging)**: 70A (for 24/3000/70 model), up to 120A.
* **Capacity/Power**: 800VA up to 5000VA (e.g., 2400W continuous for 3000VA model).
* **Key Ports/Connections**: DC battery terminals, AC-in, AC-out-1 (no-break), AC-out-2, VE.Bus.
* **Diagram-Friendly Image**: [Victron MultiPlus](https://www.victronenergy.com/upload/products/MultiPlus-12-3000-120-front.png)

### Renogy Pure Sine Wave Inverter
* **Description**: Standalone pure sine wave inverter (UPS models available).
* **Voltage**: 12V DC input, 110V/120V AC output.
* **Max Amps**: N/A (Discharge depends on load, e.g., ~166A for 2000W at 12V).
* **Capacity/Power**: 1000W, 2000W, or 3000W continuous.
* **Key Ports/Connections**: DC input terminals, AC outlets, hardwire terminal block, remote port.
* **Diagram-Friendly Image**: [Renogy Inverter](https://www.renogy.com/content/images/Inverter/2000W-12V-Pure-Sine-Wave-Inverter/2000W-12V-Pure-Sine-Wave-Inverter-1.jpg)

---

## 2. Alternator

### Generic Alternator
* **Description**: Standard marine high-output alternator.
* **Voltage**: 12V or 24V.
* **Max Amps**: Typically 100A to 400A (Requires an `amp_rating` field).
* **Key Ports/Connections**: B+ (Positive output), B- (Ground), Field (Excitation), Stator/Tachometer.
* **Diagram-Friendly Image**: [Generic Alternator](https://www.mastervolt.com/images/products/AlphaCompact.jpg)

---

## 3. Alternator Regulators

### Wakespeed WS500 Pro
* **Description**: Advanced alternator regulator with Bluetooth and CANbus.
* **Voltage**: Auto-detects 12V, 24V, and 48V systems.
* **Max Amps**: Drives field current for high-output alternators (up to 500A+ shunt sensing).
* **Key Ports/Connections**: Ampseal 23-pin waterproof connector, RJ45 (CANbus), USB Type B.
* **Diagram-Friendly Image**: [Wakespeed WS500 Pro](https://www.wakespeed.com/wp-content/uploads/2023/11/WS500-PRO-1.png)

### Mastervolt Alpha Pro III
* **Description**: 3-step charge regulator for standard and high-performance alternators.
* **Voltage**: 12V and 24V selectable.
* **Max Amps**: 20A field current capability (compatible with alternators up to 400A).
* **Key Ports/Connections**: Alpha/Mastervolt connection plug, MasterBus, Temp sensor port.
* **Diagram-Friendly Image**: [Mastervolt Alpha Pro III](https://www.mastervolt.com/images/products/AlphaProIII.jpg)

---

## 4. Batteries

### Fogstar Drift Gen 2 628Ah
* **Description**: High-capacity LiFePO4 leisure battery with active balancing and heating.
* **Voltage**: 12.8V nominal.
* **Max Amps**: 300A max continuous discharge.
* **Capacity**: 628Ah (8,038Wh).
* **Key Ports/Connections**: M8 terminals (1.25 pitch).
* **Diagram-Friendly Image**: [Fogstar Drift Gen 2](https://www.fogstar.co.uk/cdn/shop/products/Drift628Ah_1024x1024.jpg)

---

## 5. MPPT Solar Charge Controllers

### Victron SmartSolar MPPT (e.g., 100/50)
* **Description**: Solar charge controller with ultra-fast Maximum Power Point Tracking.
* **Voltage**: 12V/24V auto-select battery voltage; up to 100V PV input.
* **Max Amps**: 50A rated charge current.
* **Key Ports/Connections**: PV+ / PV- terminals, Batt+ / Batt- terminals, VE.Direct port.
* **Diagram-Friendly Image**: [Victron SmartSolar MPPT](https://www.victronenergy.com/upload/products/SmartSolar%20MPPT%20100-50%20(top).png)

---

## 6. System Monitor

### Victron Cerbo GX MK2
* **Description**: Central communication and monitoring center for Victron systems.
* **Voltage**: 8-70V DC supply voltage.
* **Max Amps**: Low power consumption (2.8W at 12V).
* **Key Ports/Connections**: 3x VE.Direct, 2x VE.Bus, 2x VE.Can, Ethernet, 3x USB, HDMI, tank/temp inputs, relays.
* **Diagram-Friendly Image**: [Victron Cerbo GX MK2](https://cdn11.bigcommerce.com/s-6rtev5owwo/images/stencil/1280x1280/products/441707/369350/102424XL__79040.1739477034.jpg)

---

## 7. Starlink 12V Power Supply

### Starlink 12V DC Conversion Kit (e.g., Trio Flatmount / Starvmount)
* **Description**: DC-to-DC step-up converter to power Starlink directly from a 12V battery bank, bypassing the inverter.
* **Voltage**: 12V-24V DC input, 56V DC output.
* **Max Amps**: ~3A to 5A output (requires 10A-15A fusing on the 12V side).
* **Capacity/Power**: ~168W to 336W.
* **Key Ports/Connections**: 12V DC input (+/-), RJ45/PoE output to Starlink dish/router.
* **Diagram-Friendly Image**: [Starlink 12V Conversion Kit](https://www.trioflatmount.com/cdn/shop/files/Gen312V.png)

---

## 8. Shunt

### Victron SmartShunt
* **Description**: Intelligent battery monitor without a display (uses Bluetooth/VE.Direct).
* **Voltage**: 6.5 - 70 VDC supply voltage.
* **Max Amps**: 500A (also available in 300A, 1000A, 2000A).
* **Key Ports/Connections**: M10 bolts for battery/system minus, VE.Direct port, Aux input (temp/midpoint).
* **Diagram-Friendly Image**: [Victron SmartShunt](https://www.victronenergy.com/upload/images/SmartShunt-500A-front.png)

---

## JSON Schema Mapping

To integrate these real-world components into the app's schematic engine, their specifications map to the following JSON schema fields:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "id": { "type": "string", "description": "Unique identifier (e.g., 'victron_cerbo_gx_mk2')" },
    "category": { "type": "string", "enum": ["inverter", "alternator", "regulator", "battery", "mppt", "monitor", "power_supply", "shunt"] },
    "manufacturer": { "type": "string", "description": "e.g., 'Victron', 'Renogy', 'Fogstar'" },
    "model": { "type": "string", "description": "e.g., 'SmartSolar MPPT 100/50'" },
    "image_url": { "type": "string", "format": "uri", "description": "URL to top-down diagram-friendly image" },
    
    "electrical_specs": {
      "type": "object",
      "properties": {
        "nominal_voltage_dc": { "type": "number", "description": "e.g., 12, 24, 48" },
        "operating_voltage_range": { 
          "type": "array", 
          "items": { "type": "number" },
          "description": "[min_voltage, max_voltage], e.g., [8, 70] for Cerbo GX"
        },
        "max_continuous_amps": { "type": "number", "description": "e.g., 300 for Fogstar Drift" },
        "capacity_ah": { "type": "number", "description": "e.g., 628 for Fogstar Drift" },
        "power_watts": { "type": "number", "description": "e.g., 2500 for Sterling Inverter" },
        "amp_rating": { "type": "number", "description": "Specific to Alternators" }
      }
    },
    
    "ports": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "port_id": { "type": "string" },
          "port_type": { "type": "string", "enum": ["dc_positive", "dc_negative", "ac_in", "ac_out", "ve_direct", "ve_bus", "canbus", "rj45", "custom"] },
          "max_amps": { "type": "number", "description": "Max current for this specific port" }
        }
      }
    }
  }
}
```

### Mapping Examples
* **Fogstar Drift Gen 2 628Ah**: Maps `capacity_ah` to `628`, `nominal_voltage_dc` to `12.8`, and `max_continuous_amps` to `300`. Ports include `dc_positive` and `dc_negative` (M8 terminals).
* **Victron SmartSolar MPPT 100/50**: Maps `operating_voltage_range` to `[12, 24]`, `max_continuous_amps` (charge) to `50`. Ports include `dc_positive` (PV), `dc_negative` (PV), `dc_positive` (Batt), `dc_negative` (Batt), and `ve_direct`.
* **Generic Alternator**: Maps `amp_rating` to the user-defined output (e.g., `150`). Ports include `dc_positive` (B+), `dc_negative` (B-), and `custom` (Field/Stator).
