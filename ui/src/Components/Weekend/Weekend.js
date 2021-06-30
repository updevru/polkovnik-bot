import React from 'react';
import {Form, Button, DatePicker, Space, Checkbox} from 'antd';
import Dictionary from "../../Services/Dictionary";
import {MinusCircleOutlined, PlusOutlined} from "@ant-design/icons";
import moment from "moment";

const { RangePicker } = DatePicker;

export function WeekendFormValue (result)
{
    if (result && 'weekend' in result && typeof result.weekend.intervals === "object" && result.weekend.intervals) {
        result.weekend.intervals = result.weekend.intervals.map(function (row){
            if ('start' in row) {
                return [moment(row.start, 'DD-MM-YYYY'), moment(row.end, 'DD-MM-YYYY')]
            }
            return row;
        });
    }

    return result
}

export function WeekendDataValue(values)
{
    if (typeof values.weekend.intervals === "object" && values.weekend.intervals) {
        values.weekend.intervals = values.weekend.intervals.map(function (row){
            return {
                start: row[0].format("DD-MM-YYYY"),
                end: row[1].format("DD-MM-YYYY")
            }
        });
    }

    return values;
}

export class Weekend extends React.Component {

    static defaultProps = {
        button_add_title: 'Добавить отпуск'
    }

    render() {
        return <span>
            <Form.Item
                label="Не рабочие дни"
                name={['weekend', 'week_days']}
            >
                    <Checkbox.Group options={Dictionary.getWeekdays()} />
                </Form.Item>

                <Form.List name={['weekend', 'intervals']}>
                    {(fields, { add, remove }) => (
                        <>
                            {fields.map(({ key, name, fieldKey, ...restField }) => (
                                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                                    <Form.Item
                                        {...restField}
                                        name={name}
                                        fieldKey={name}
                                        rules={[{ required: true, message: 'Missing interval' }]}
                                    >
                                        <RangePicker format={"DD-MM-YYYY"} />
                                    </Form.Item>
                                    <MinusCircleOutlined onClick={() => remove(name)} />
                                </Space>
                            ))}
                            <Form.Item>
                                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                                    {this.props.button_add_title}
                                </Button>
                            </Form.Item>
                        </>
                    )}
                </Form.List>
        </span>
    }
}
